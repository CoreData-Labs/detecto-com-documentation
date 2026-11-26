// Command main downloads every PDF document that is reachable from a fixed list of
// Detecto product pages and stores the files, with sanitized names, in a "PDFs" directory.
//
// Run it with:  go run main.go
//
// Only the Go standard library is used, so no "go get" or module setup is required
// beyond what "go run" already handles for a single file.
package main // Every executable Go program lives in package main.

import ( // Start of the import block; each line imports one standard library package.
	"errors"        // Creating, wrapping and inspecting errors (errors.New, errors.As).
	"fmt"           // Formatting text and building error messages.
	"html"          // Converting HTML entities such as &amp; back into normal characters.
	"io"            // Streaming data from HTTP responses into files.
	"log"           // Console logging with timestamps.
	"net/http"      // Sending HTTP requests.
	"net/url"       // Parsing, resolving and normalizing URLs.
	"os"            // Creating directories, files, and exiting with a status code.
	"path"          // Extracting the last segment of a URL path (always uses forward slashes).
	"path/filepath" // Building operating-system specific file paths.
	"regexp"        // Finding links inside HTML text (the standard library has no HTML parser).
	"strconv"       // Converting the Retry-After header text into a number.
	"strings"       // Text manipulation helpers.
	"time"          // Timeouts, delays between requests and retry back-off.
) // End of the import block.

// ---------------------------------------------------------------------------------------
// Configuration constants
// ---------------------------------------------------------------------------------------

const outputDirectoryName = "PDFs" // Name of the directory (inside the working directory) that receives every PDF.

const userAgentText = "Mozilla/5.0 (X11; CrOS x86_64 14541.0.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36" // Browser-style User-Agent requested for every HTTP request.

const httpRequestTimeout = 90 * time.Second // Upper limit for one complete HTTP request, including reading the body.

const maximumAttemptsPerRequest = 4 // How many times one request is tried before it is reported as failed.

const initialRetryDelay = 2 * time.Second // Waiting time before the first retry; it doubles on every further retry.

const maximumRetryDelay = 30 * time.Second // Upper limit for any single waiting period between retries.

const delayBetweenRequests = 750 * time.Millisecond // Minimum pause between two requests so the server is not flooded.

const maximumPaginationPagesPerCategory = 20 // Safety limit that prevents endless pagination loops.

const maximumProductPagesPerCategory = 25 // Safety limit for how many product pages are visited per category.

const maximumPageBodyBytes = 10 * 1024 * 1024 // HTML pages larger than 10 MiB are truncated to protect memory.

const pdfSignatureText = "%PDF-" // Every valid PDF file contains this marker near its beginning.

const pdfSignatureSearchLength = 1024 // The PDF specification allows the marker anywhere within the first 1024 bytes.

const fallbackFilenameBase = "document" // Base name used when sanitizing leaves nothing usable.

// categoryPageAddresses lists every page that the program must visit.
var categoryPageAddresses = []string{ // A slice (list) of text values.
	"https://detecto.com/category/baby-scales/subcategory/6745-digital",
	"https://detecto.com/category/baby-scales/subcategory/mb130-digital",
	"https://detecto.com/category/baby-scales/subcategory/8440-digital",
	"https://detecto.com/category/baby-scales/subcategory/8435-digital",
	"https://detecto.com/category/baby-scales/subcategory/450-series-weigh-beam",
	"https://detecto.com/category/baby-scales/subcategory/hs25kgp-suspended",
} // End of the list of category pages.

// productPathFragments identifies links that lead to individual product detail pages.
// Generic navigation links such as "/products/baby-scales" intentionally do NOT match,
// otherwise the crawler would wander through the whole website menu.
var productPathFragments = []string{"/product/", "/subcategory/", "/resource/"} // Adjust this list if the website structure changes.

// ---------------------------------------------------------------------------------------
// Regular expressions used to find links inside HTML (compiled once at program start)
// ---------------------------------------------------------------------------------------

// htmlAttributeValuePattern captures the value of href, src and data-* attributes, quoted with " or '.
var htmlAttributeValuePattern = regexp.MustCompile(`(?i)\b(?:href|src|data-[a-z0-9_-]+)\s*=\s*(?:"([^"]*)"|'([^']*)')`)

// embeddedAbsolutePdfPattern finds absolute PDF URLs anywhere in the text, e.g. inside JSON
// that modern JavaScript websites (such as Next.js) embed in the HTML page.
var embeddedAbsolutePdfPattern = regexp.MustCompile(`(?i)https?://[^\s"'<>\\]+?\.pdf(?:\?[^\s"'<>\\]*)?`)

// embeddedRootRelativePdfPattern finds quoted root-relative PDF paths such as "/files/manual.pdf".
var embeddedRootRelativePdfPattern = regexp.MustCompile(`(?i)["'](/[^"'\s<>\\]*\.pdf)["'\\]`)

// jsonEscapeReplacer undoes common JSON string escapes so that embedded URLs become ordinary URLs.
var jsonEscapeReplacer = strings.NewReplacer(`\/`, `/`, `\u0026`, `&`, `\u002f`, `/`, `\u002F`, `/`)

// ---------------------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------------------

// runStatistics collects the numbers that are printed in the final summary.
type runStatistics struct {
	categoryPagesProcessed  int // Category pages (including extra pagination pages) fetched successfully.
	productPagesVisited     int // Individual product pages fetched successfully.
	pageFetchFailures       int // Pages that could not be fetched at all.
	uniquePdfUrlsDiscovered int // Distinct PDF URLs found across all pages.
	pdfsDownloaded          int // PDFs saved successfully during this run.
	pdfsSkipped             int // PDFs not downloaded because the target file already exists.
	pdfsFailed              int // PDFs whose download or validation failed.
}

// temporaryError marks a failure that is worth retrying (network problem, HTTP 5xx, HTTP 429, ...).
type temporaryError struct {
	cause          error         // The underlying problem, kept for logging.
	suggestedDelay time.Duration // Delay requested by the server through Retry-After, or zero.
}

// Error returns the message of the underlying problem so temporaryError satisfies the error interface.
func (problem *temporaryError) Error() string {
	// Delegate to the wrapped error so log messages stay readable.
	return problem.cause.Error()
}

// Unwrap lets errors.Is and errors.As look inside the temporaryError.
func (problem *temporaryError) Unwrap() error {
	// Expose the wrapped error.
	return problem.cause
}

// pdfDownloader bundles the HTTP client and all state that is shared while the program runs.
type pdfDownloader struct {
	httpClient                  *http.Client    // Shared HTTP client configured with a timeout.
	outputDirectoryAbsolutePath string          // Absolute path of the PDFs directory.
	statistics                  runStatistics   // Counters for the final summary.
	discoveredPdfUrls           map[string]bool // Normalized PDF URLs already discovered (prevents duplicate downloads).
	visitedPageUrls             map[string]bool // Normalized page URLs already fetched (prevents loops).
	filenamesClaimedThisRun     map[string]bool // Final filenames already assigned to a PDF during this run.
	timeOfLastRequest           time.Time       // Used to enforce delayBetweenRequests.
}

// ---------------------------------------------------------------------------------------
// Logging helper
// ---------------------------------------------------------------------------------------

// writeDebugLog prints one debug-level message. The standard log package has no levels,
// so the "[DEBUG]" prefix is added manually.
func writeDebugLog(messageFormat string, messageArguments ...any) {
	// Prepend the level label and print using the standard logger (stderr, with timestamp).
	log.Printf("[DEBUG] "+messageFormat, messageArguments...)
}

// ---------------------------------------------------------------------------------------
// Filename and URL helpers
// ---------------------------------------------------------------------------------------

// sanitizeFilename converts any filename into a safe one that contains only a-z, 0-9 and
// single underscores, never starts or ends with an underscore, and always ends in ".pdf".
func sanitizeFilename(originalFilename string) string {
	// Lowercase everything and remove surrounding whitespace.
	lowercaseName := strings.ToLower(strings.TrimSpace(originalFilename))
	// Remove the extension so it can be re-added in one consistent form at the end.
	lowercaseName = strings.TrimSuffix(lowercaseName, ".pdf")
	// A builder efficiently assembles the cleaned text.
	var cleanedNameBuilder strings.Builder
	// Starting with "true" prevents a leading underscore (invalid leading characters are simply dropped).
	previousCharacterWasUnderscore := true
	// Inspect the name one character at a time.
	for _, character := range lowercaseName {
		// Letters a-z and digits 0-9 are the only characters kept as they are.
		if (character >= 'a' && character <= 'z') || (character >= '0' && character <= '9') {
			// Keep the allowed character.
			cleanedNameBuilder.WriteRune(character)
			// Remember that the last written character was not an underscore.
			previousCharacterWasUnderscore = false
		} else if !previousCharacterWasUnderscore {
			// Replace a run of invalid characters with exactly one underscore (prevents "__").
			cleanedNameBuilder.WriteRune('_')
			// Remember that an underscore was just written.
			previousCharacterWasUnderscore = true
		}
	}
	// Remove a trailing underscore, which may remain when the name ended with invalid characters.
	cleanedNameBase := strings.TrimRight(cleanedNameBuilder.String(), "_")
	// Edge case: nothing usable was left (for example the original name was only symbols).
	if cleanedNameBase == "" {
		// Use a safe fallback; uniqueness is guaranteed later by reserveUniqueFilename.
		cleanedNameBase = fallbackFilenameBase
	}
	// Re-attach the lowercase .pdf extension.
	return cleanedNameBase + ".pdf"
}

// normalizeURL resolves a possibly relative link against the page it was found on and
// returns a canonical absolute URL string (no fragment, lowercase host, http/https only).
func normalizeURL(rawLink string, baseURL *url.URL) (string, error) {
	// Parse the raw link text into a URL structure.
	parsedReference, parseError := url.Parse(rawLink)
	// A link that cannot be parsed is reported to the caller.
	if parseError != nil {
		// Wrap the error with context.
		return "", fmt.Errorf("cannot parse link %q: %w", rawLink, parseError)
	}
	// Turn relative links ("/files/a.pdf", "../a.pdf") into absolute ones using the page URL.
	resolvedURL := baseURL.ResolveReference(parsedReference)
	// Only web links are interesting; mailto:, tel:, javascript: and similar are rejected.
	if resolvedURL.Scheme != "http" && resolvedURL.Scheme != "https" {
		// Report why the link was rejected.
		return "", fmt.Errorf("unsupported URL scheme %q in link %q", resolvedURL.Scheme, rawLink)
	}
	// A link without a host cannot be requested.
	if resolvedURL.Host == "" {
		// Report why the link was rejected.
		return "", fmt.Errorf("link %q has no host name", rawLink)
	}
	// Fragments (#section) never change the downloaded document, so drop them.
	resolvedURL.Fragment = ""
	// Host names are case-insensitive; lowercase them so duplicates compare equal.
	resolvedURL.Host = strings.ToLower(resolvedURL.Host)
	// An empty path means the site root.
	if resolvedURL.Path == "" {
		// Use "/" so equal addresses have one single textual form.
		resolvedURL.Path = "/"
	}
	// Return the canonical text form.
	return resolvedURL.String(), nil
}

// isPdfURL reports whether the path of the URL ends with ".pdf" (query strings are ignored).
func isPdfURL(parsedURL *url.URL) bool {
	// Compare case-insensitively because some servers use ".PDF".
	return strings.HasSuffix(strings.ToLower(parsedURL.Path), ".pdf")
}

// deriveOriginalFilename extracts the last path segment of a PDF URL, e.g. "Baby Scale 6745.pdf".
func deriveOriginalFilename(pdfAddress string) string {
	// Parse the address; failures yield an empty name, which sanitizeFilename turns into a fallback.
	parsedURL, parseError := url.Parse(pdfAddress)
	// Handle the unlikely parse failure explicitly instead of ignoring it.
	if parseError != nil {
		// Log the problem and fall back to an empty name.
		writeDebugLog("Could not parse %q to derive a filename: %v", pdfAddress, parseError)
		// An empty result triggers the safe fallback name later.
		return ""
	}
	// parsedURL.Path is already percent-decoded ("%20" became a space); take its last segment.
	lastSegment := path.Base(parsedURL.Path)
	// path.Base returns "." or "/" when there is no real file name.
	if lastSegment == "." || lastSegment == "/" {
		// Treat this as "no name available".
		return ""
	}
	// Return the raw last segment.
	return lastSegment
}

// extractCandidateLinks finds every link-like value in an HTML document and returns the
// normalized, de-duplicated absolute URLs in the order they first appear.
func extractCandidateLinks(htmlText string, pageURL *url.URL) []string {
	// Undo JSON escapes so URLs hidden in embedded script data are readable.
	unescapedText := jsonEscapeReplacer.Replace(htmlText)
	// Collect raw link values from several sources.
	var rawLinkValues []string
	// Source 1: the values of href/src/data-* attributes.
	for _, attributeMatch := range htmlAttributeValuePattern.FindAllStringSubmatch(unescapedText, -1) {
		// Group 1 holds double-quoted values and group 2 holds single-quoted values.
		attributeValue := attributeMatch[1]
		// If the double-quoted group is empty, use the single-quoted group.
		if attributeValue == "" {
			// Take the single-quoted value.
			attributeValue = attributeMatch[2]
		}
		// Remember the value.
		rawLinkValues = append(rawLinkValues, attributeValue)
	}
	// Source 2: absolute PDF URLs found anywhere in the text (covers JSON embedded in scripts).
	rawLinkValues = append(rawLinkValues, embeddedAbsolutePdfPattern.FindAllString(unescapedText, -1)...)
	// Source 3: quoted root-relative PDF paths found anywhere in the text.
	for _, relativeMatch := range embeddedRootRelativePdfPattern.FindAllStringSubmatch(unescapedText, -1) {
		// Group 1 is the path without the surrounding quote characters.
		rawLinkValues = append(rawLinkValues, relativeMatch[1])
	}
	// Track which URLs were already added so the result has no duplicates.
	alreadyAdded := make(map[string]bool)
	// This slice receives the final answer.
	var normalizedLinks []string
	// Process each raw value.
	for _, rawValue := range rawLinkValues {
		// Convert HTML entities (&amp;) and trim whitespace.
		cleanedValue := strings.TrimSpace(html.UnescapeString(rawValue))
		// Skip values that are empty or obviously not web page links.
		if cleanedValue == "" || strings.HasPrefix(cleanedValue, "#") ||
			strings.HasPrefix(cleanedValue, "mailto:") || strings.HasPrefix(cleanedValue, "tel:") ||
			strings.HasPrefix(cleanedValue, "javascript:") || strings.HasPrefix(cleanedValue, "data:") {
			// Move on to the next value.
			continue
		}
		// Resolve and normalize the link.
		normalizedLink, normalizeError := normalizeURL(cleanedValue, pageURL)
		// Links that cannot be normalized are logged (not silently dropped) and skipped.
		if normalizeError != nil {
			// Debug-level message because many attribute values are legitimately not URLs.
			writeDebugLog("Skipping unusable link value: %v", normalizeError)
			// Move on to the next value.
			continue
		}
		// Ignore links already collected.
		if alreadyAdded[normalizedLink] {
			// Move on to the next value.
			continue
		}
		// Remember that this link has been added.
		alreadyAdded[normalizedLink] = true
		// Append it to the result list.
		normalizedLinks = append(normalizedLinks, normalizedLink)
	}
	// Return every unique link found.
	return normalizedLinks
}

// filterPdfLinks keeps only the links that point to PDF documents (on any host, because
// Detecto stores its documents on a separate storage host).
func filterPdfLinks(allLinks []string) []string {
	// Result list.
	var pdfLinks []string
	// Examine every link.
	for _, link := range allLinks {
		// Parse again to look at the path component.
		parsedLink, parseError := url.Parse(link)
		// Links were normalized earlier, so a failure here is unexpected but must still be logged.
		if parseError != nil {
			// Log and skip.
			writeDebugLog("Unexpected parse failure for %q: %v", link, parseError)
			// Move on.
			continue
		}
		// Keep the link when it ends with .pdf.
		if isPdfURL(parsedLink) {
			// Add to the result.
			pdfLinks = append(pdfLinks, link)
		}
	}
	// Return the PDF links.
	return pdfLinks
}

// filterProductLinks keeps same-site links that look like individual product detail pages.
func filterProductLinks(allLinks []string, currentPageURL *url.URL) []string {
	// Result list.
	var productLinks []string
	// Examine every link.
	for _, link := range allLinks {
		// Parse the link to inspect host and path.
		parsedLink, parseError := url.Parse(link)
		// Skip (and log) links that cannot be parsed.
		if parseError != nil {
			// Log and skip.
			writeDebugLog("Unexpected parse failure for %q: %v", link, parseError)
			// Move on.
			continue
		}
		// Only follow links that stay on the same website.
		if !strings.EqualFold(parsedLink.Hostname(), currentPageURL.Hostname()) {
			// Different host: ignore.
			continue
		}
		// Never treat PDFs as pages.
		if isPdfURL(parsedLink) {
			// PDF: ignore here.
			continue
		}
		// Do not revisit the page we are already on.
		if parsedLink.Path == currentPageURL.Path {
			// Same page: ignore.
			continue
		}
		// Check whether the path contains one of the product-page markers.
		for _, pathFragment := range productPathFragments {
			// A match means this is probably a product detail page.
			if strings.Contains(parsedLink.Path, pathFragment) {
				// Remember the link, without its query string noise being altered.
				productLinks = append(productLinks, link)
				// One matching marker is enough.
				break
			}
		}
	}
	// Return the product links.
	return productLinks
}

// filterPaginationLinks keeps links that look like "next page" links of the current listing:
// same path with a page query parameter, or a "<path>/page/<n>" style address.
func filterPaginationLinks(allLinks []string, currentPageURL *url.URL) []string {
	// Result list.
	var paginationLinks []string
	// Examine every link.
	for _, link := range allLinks {
		// Parse the link to inspect host, path and query.
		parsedLink, parseError := url.Parse(link)
		// Skip (and log) links that cannot be parsed.
		if parseError != nil {
			// Log and skip.
			writeDebugLog("Unexpected parse failure for %q: %v", link, parseError)
			// Move on.
			continue
		}
		// Pagination always stays on the same host.
		if !strings.EqualFold(parsedLink.Hostname(), currentPageURL.Hostname()) {
			// Different host: ignore.
			continue
		}
		// Read the query parameters once.
		queryValues := parsedLink.Query()
		// Variant 1: same path plus a page-number query parameter.
		hasPageQuery := parsedLink.Path == currentPageURL.Path && (queryValues.Get("page") != "" || queryValues.Get("p") != "")
		// Variant 2: "/some/path/page/2".
		hasPagePath := strings.HasPrefix(parsedLink.Path, strings.TrimRight(currentPageURL.Path, "/")+"/page/")
		// Keep the link when either variant matches.
		if hasPageQuery || hasPagePath {
			// Add to the result.
			paginationLinks = append(paginationLinks, link)
		}
	}
	// Return the pagination links.
	return paginationLinks
}

// ---------------------------------------------------------------------------------------
// HTTP helpers
// ---------------------------------------------------------------------------------------

// closeResponseBody closes an HTTP response body and logs (instead of ignoring) any error.
func closeResponseBody(responseBody io.Closer) {
	// Closing releases the network connection back to the pool.
	if closeError := responseBody.Close(); closeError != nil {
		// Report the problem; it is not fatal for the program.
		writeDebugLog("Error while closing response body: %v", closeError)
	}
}

// waitForRequestSpacing sleeps if necessary so requests are at least delayBetweenRequests apart.
func (downloader *pdfDownloader) waitForRequestSpacing() {
	// Measure how long ago the previous request started.
	elapsedTime := time.Since(downloader.timeOfLastRequest)
	// If it was too recent, wait for the remaining time.
	if elapsedTime < delayBetweenRequests {
		// Sleep to keep the request rate polite.
		time.Sleep(delayBetweenRequests - elapsedTime)
	}
	// Record the start time of the upcoming request.
	downloader.timeOfLastRequest = time.Now()
}

// sendGetRequest performs one HTTP GET. A network failure is returned as a temporaryError.
// The caller must close the response body.
func (downloader *pdfDownloader) sendGetRequest(address string) (*http.Response, error) {
	// Log the outgoing request.
	writeDebugLog("HTTP GET %s", address)
	// Build the request object.
	request, requestError := http.NewRequest(http.MethodGet, address, nil)
	// An invalid address is a permanent problem (retrying cannot help).
	if requestError != nil {
		// Return a plain error so no retry happens.
		return nil, fmt.Errorf("cannot build request for %s: %w", address, requestError)
	}
	// Identify this program to the server.
	request.Header.Set("User-Agent", userAgentText)
	// State which content types are acceptable.
	request.Header.Set("Accept", "text/html,application/pdf,application/xhtml+xml;q=0.9,*/*;q=0.8")
	// Send the request.
	response, responseError := downloader.httpClient.Do(request)
	// Timeouts, DNS failures and dropped connections are usually temporary.
	if responseError != nil {
		// Mark as retryable.
		return nil, &temporaryError{cause: fmt.Errorf("network error for %s: %w", address, responseError)}
	}
	// Log the status code of the response.
	writeDebugLog("HTTP response for %s: %s", address, response.Status)
	// Hand the response to the caller.
	return response, nil
}

// classifyStatusCode returns nil for success, a temporaryError for statuses worth retrying
// (408, 429 and 5xx), and a permanent error for all other statuses (404, 403, ...).
func classifyStatusCode(response *http.Response) error {
	// 2xx means success.
	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		// No error.
		return nil
	}
	// Build a descriptive error message.
	statusError := fmt.Errorf("unexpected HTTP status %s for %s", response.Status, response.Request.URL)
	// Request timeout, rate limiting and server errors are considered temporary.
	if response.StatusCode == http.StatusRequestTimeout || response.StatusCode == http.StatusTooManyRequests || response.StatusCode >= 500 {
		// Honor a Retry-After header expressed in seconds, if the server sent one.
		var serverRequestedDelay time.Duration
		// Read the header text.
		retryAfterText := response.Header.Get("Retry-After")
		// Convert it to a number when possible.
		if retryAfterSeconds, conversionError := strconv.Atoi(retryAfterText); conversionError == nil && retryAfterSeconds > 0 {
			// Store the delay requested by the server.
			serverRequestedDelay = time.Duration(retryAfterSeconds) * time.Second
		}
		// Mark as retryable.
		return &temporaryError{cause: statusError, suggestedDelay: serverRequestedDelay}
	}
	// Everything else (such as 404 or 403) will not improve by retrying.
	return statusError
}

// retryOperation runs an operation up to maximumAttemptsPerRequest times, retrying only
// when it fails with a temporaryError, and using exponential back-off between attempts.
func (downloader *pdfDownloader) retryOperation(operationDescription string, operation func() error) error {
	// Remember the most recent failure so it can be reported.
	var lastError error
	// Try the operation repeatedly.
	for attemptNumber := 1; attemptNumber <= maximumAttemptsPerRequest; attemptNumber++ {
		// Respect the polite request spacing before every attempt.
		downloader.waitForRequestSpacing()
		// Run the real work.
		lastError = operation()
		// Success ends the loop immediately.
		if lastError == nil {
			// Report success.
			return nil
		}
		// Determine whether the failure is temporary.
		var temporaryProblem *temporaryError
		// errors.As unwraps nested errors to find a temporaryError.
		if !errors.As(lastError, &temporaryProblem) {
			// Permanent failure: do not retry.
			writeDebugLog("%s failed permanently (attempt %d): %v", operationDescription, attemptNumber, lastError)
			// Give up immediately.
			return lastError
		}
		// No waiting is needed after the final attempt.
		if attemptNumber == maximumAttemptsPerRequest {
			// Leave the loop and report the failure below.
			break
		}
		// Exponential back-off: 2s, 4s, 8s, ...
		retryDelay := initialRetryDelay * time.Duration(1<<uint(attemptNumber-1))
		// Prefer a longer delay if the server asked for one.
		if temporaryProblem.suggestedDelay > retryDelay {
			// Use the server's delay.
			retryDelay = temporaryProblem.suggestedDelay
		}
		// Never wait longer than the configured maximum.
		if retryDelay > maximumRetryDelay {
			// Cap the delay.
			retryDelay = maximumRetryDelay
		}
		// Log the retry decision.
		writeDebugLog("%s failed (attempt %d of %d): %v. Retrying in %s", operationDescription, attemptNumber, maximumAttemptsPerRequest, lastError, retryDelay)
		// Wait before the next attempt.
		time.Sleep(retryDelay)
	}
	// All attempts were used.
	return fmt.Errorf("%s failed after %d attempts: %w", operationDescription, maximumAttemptsPerRequest, lastError)
}

// fetchPageBody downloads an HTML page (with retries) and returns its text.
func (downloader *pdfDownloader) fetchPageBody(pageAddress string) (string, error) {
	// The closure below stores the page text here.
	var pageBody string
	// Run the request through the retry helper.
	fetchError := downloader.retryOperation("fetch page "+pageAddress, func() error {
		// Send the request.
		response, requestError := downloader.sendGetRequest(pageAddress)
		// Pass failures to the retry helper.
		if requestError != nil {
			// Return the (possibly temporary) error.
			return requestError
		}
		// Always close the body when this attempt ends.
		defer closeResponseBody(response.Body)
		// Check the HTTP status code.
		if statusError := classifyStatusCode(response); statusError != nil {
			// Return the (possibly temporary) error.
			return statusError
		}
		// Read the body, limited in size to protect memory.
		bodyBytes, readError := io.ReadAll(io.LimitReader(response.Body, maximumPageBodyBytes))
		// A read failure usually means an interrupted connection, so it is retryable.
		if readError != nil {
			// Mark as retryable.
			return &temporaryError{cause: fmt.Errorf("could not read body of %s: %w", pageAddress, readError)}
		}
		// Convert the bytes to text for link extraction.
		pageBody = string(bodyBytes)
		// Success.
		return nil
	})
	// Return the text and any error.
	return pageBody, fetchError
}

// ---------------------------------------------------------------------------------------
// File handling
// ---------------------------------------------------------------------------------------

// ensureOutputDirectory creates the PDFs directory when needed and returns its absolute path.
func ensureOutputDirectory(directoryName string) (string, error) {
	// Convert the relative name into an absolute path based on the working directory.
	absolutePath, absoluteError := filepath.Abs(directoryName)
	// Report a failure to determine the working directory.
	if absoluteError != nil {
		// Wrap with context.
		return "", fmt.Errorf("cannot determine absolute path of %q: %w", directoryName, absoluteError)
	}
	// Check whether the directory already exists.
	if _, statError := os.Stat(absolutePath); statError == nil {
		// Nothing to create.
		writeDebugLog("Output directory already exists: %s", absolutePath)
		// Return the path.
		return absolutePath, nil
	} else if !os.IsNotExist(statError) {
		// Any error other than "does not exist" (e.g. permission denied) is a real problem.
		return "", fmt.Errorf("cannot inspect directory %s: %w", absolutePath, statError)
	}
	// Create the directory (and parents) with normal permissions.
	if createError := os.MkdirAll(absolutePath, 0o755); createError != nil {
		// Wrap with context.
		return "", fmt.Errorf("cannot create directory %s: %w", absolutePath, createError)
	}
	// Log the creation.
	writeDebugLog("Created output directory: %s", absolutePath)
	// Return the path.
	return absolutePath, nil
}

// reserveUniqueFilename returns a filename that no other PDF in this run has used, adding
// _2, _3, ... before ".pdf" when two different PDFs sanitize to the same name.
func (downloader *pdfDownloader) reserveUniqueFilename(sanitizedFilename string) string {
	// Start with the sanitized name itself.
	candidateFilename := sanitizedFilename
	// Separate the base name from the extension so a counter can be inserted.
	baseName := strings.TrimSuffix(sanitizedFilename, ".pdf")
	// Counter used for suffixes.
	duplicateCounter := 1
	// Keep trying until an unused name is found.
	for downloader.filenamesClaimedThisRun[candidateFilename] {
		// Move to the next suffix number.
		duplicateCounter++
		// Build e.g. "baby_scale_2.pdf" (the suffix keeps all filename rules intact).
		candidateFilename = baseName + "_" + strconv.Itoa(duplicateCounter) + ".pdf"
		// Log the collision resolution.
		writeDebugLog("Duplicate filename detected for %q; trying %q", sanitizedFilename, candidateFilename)
	}
	// Claim the name so later PDFs cannot take it.
	downloader.filenamesClaimedThisRun[candidateFilename] = true
	// Return the unique name.
	return candidateFilename
}

// validatePdfFile checks that the file is non-empty and contains the PDF signature near its start.
func validatePdfFile(filePath string) error {
	// Open the file for reading.
	file, openError := os.Open(filePath)
	// Report open failures.
	if openError != nil {
		// Wrap with context.
		return fmt.Errorf("cannot open %s for validation: %w", filePath, openError)
	}
	// Close the file when validation ends; report close problems.
	defer func() {
		// Closing a read-only file rarely fails, but it is still reported.
		if closeError := file.Close(); closeError != nil {
			// Log the problem.
			writeDebugLog("Error closing %s after validation: %v", filePath, closeError)
		}
	}()
	// Look up the file size.
	fileInformation, statError := file.Stat()
	// Report stat failures.
	if statError != nil {
		// Wrap with context.
		return fmt.Errorf("cannot read file information of %s: %w", filePath, statError)
	}
	// An empty file can never be a valid PDF.
	if fileInformation.Size() == 0 {
		// Report the problem.
		return fmt.Errorf("file %s is empty", filePath)
	}
	// Prepare a buffer for the beginning of the file.
	headerBuffer := make([]byte, pdfSignatureSearchLength)
	// Read up to 1024 bytes (io.ReadFull would fail on small files, so Read is used).
	bytesRead, readError := file.Read(headerBuffer)
	// io.EOF is acceptable only when some bytes were read; anything else is an error.
	if readError != nil && !errors.Is(readError, io.EOF) {
		// Wrap with context.
		return fmt.Errorf("cannot read header of %s: %w", filePath, readError)
	}
	// Look for "%PDF-" inside the bytes that were read.
	if !strings.Contains(string(headerBuffer[:bytesRead]), pdfSignatureText) {
		// Typical cause: the server returned an HTML error page instead of a PDF.
		return fmt.Errorf("file %s does not contain a PDF signature (probably not a PDF)", filePath)
	}
	// The file looks like a valid PDF.
	return nil
}

// removeTemporaryFile deletes a leftover ".part" file; a missing file is not an error.
func removeTemporaryFile(temporaryFilePath string) {
	// Try to delete the file.
	if removeError := os.Remove(temporaryFilePath); removeError != nil && !os.IsNotExist(removeError) {
		// Report real deletion problems (a missing file simply means there is nothing to clean).
		writeDebugLog("Could not remove temporary file %s: %v", temporaryFilePath, removeError)
	}
}

// downloadOnce performs a single download attempt into a temporary file, validates it,
// and renames it to its final name. Writing to ".part" first guarantees that an
// interrupted download never leaves a half-written file under the final name.
func (downloader *pdfDownloader) downloadOnce(pdfAddress string, temporaryFilePath string, finalFilePath string) error {
	// Send the request.
	response, requestError := downloader.sendGetRequest(pdfAddress)
	// Pass failures to the retry helper.
	if requestError != nil {
		// Return the (possibly temporary) error.
		return requestError
	}
	// Always close the response body when this attempt ends.
	defer closeResponseBody(response.Body)
	// Check the HTTP status code.
	if statusError := classifyStatusCode(response); statusError != nil {
		// Return the (possibly temporary) error.
		return statusError
	}
	// Read the Content-Type header the server announced for this file.
	announcedContentType := strings.ToLower(response.Header.Get("Content-Type"))
	// Web pages, plain text and JSON answers (typical error pages) are never PDFs, so stop before writing anything.
	if strings.HasPrefix(announcedContentType, "text/") || strings.Contains(announcedContentType, "html") || strings.Contains(announcedContentType, "json") {
		// Permanent error: retrying will not turn a web page into a PDF.
		return fmt.Errorf("server announced content type %q for %s, which is not a PDF", announcedContentType, pdfAddress)
	}
	// Tracks success so the deferred cleanup knows whether to delete the temporary file.
	downloadSucceeded := false
	// Remove the temporary file on every failure path.
	defer func() {
		// Only clean up when the download did not complete.
		if !downloadSucceeded {
			// Delete the partial file.
			removeTemporaryFile(temporaryFilePath)
		}
	}()
	// Create (or truncate a stale) temporary file.
	temporaryFile, createError := os.OpenFile(temporaryFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	// File system problems are permanent for this attempt.
	if createError != nil {
		// Wrap with context.
		return fmt.Errorf("cannot create %s: %w", temporaryFilePath, createError)
	}
	// Stream the response body to disk without loading it all into memory.
	bytesWritten, copyError := io.Copy(temporaryFile, response.Body)
	// Close the file before validating or renaming it.
	closeError := temporaryFile.Close()
	// A copy error normally means the connection broke mid-download, so retry.
	if copyError != nil {
		// Mark as retryable.
		return &temporaryError{cause: fmt.Errorf("interrupted download of %s: %w", pdfAddress, copyError)}
	}
	// A close error means data may not have reached the disk.
	if closeError != nil {
		// Wrap with context.
		return fmt.Errorf("cannot finish writing %s: %w", temporaryFilePath, closeError)
	}
	// If the server announced a size, make sure we received exactly that many bytes.
	if response.ContentLength >= 0 && bytesWritten != response.ContentLength {
		// Incomplete download: retry.
		return &temporaryError{cause: fmt.Errorf("incomplete download of %s: got %d of %d bytes", pdfAddress, bytesWritten, response.ContentLength)}
	}
	// Verify the file is a real, non-empty PDF.
	if validationError := validatePdfFile(temporaryFilePath); validationError != nil {
		// Invalid content will not improve by retrying.
		return validationError
	}
	// Move the validated file to its final name.
	if renameError := os.Rename(temporaryFilePath, finalFilePath); renameError != nil {
		// Wrap with context.
		return fmt.Errorf("cannot rename %s to %s: %w", temporaryFilePath, finalFilePath, renameError)
	}
	// Tell the deferred cleanup that nothing needs to be removed.
	downloadSucceeded = true
	// Log the size of the saved file.
	writeDebugLog("Saved %d bytes to %s", bytesWritten, finalFilePath)
	// Success.
	return nil
}

// downloadPdf decides the final filename for one PDF URL and downloads it unless a file
// with that name already exists (existing files are never overwritten).
func (downloader *pdfDownloader) downloadPdf(pdfAddress string) {
	// Log the start of this download.
	writeDebugLog("Download starting: %s", pdfAddress)
	// Parse the address so its path can be checked for a .pdf extension before anything is requested.
	parsedPdfURL, parseError := url.Parse(pdfAddress)
	// Only addresses that parse correctly AND end in .pdf (any letter case) may be downloaded.
	if parseError != nil || !isPdfURL(parsedPdfURL) {
		// Explain the decision in the log (parseError is nil when the address simply lacks .pdf).
		writeDebugLog("Skipping %s: the address does not end in .pdf (parse error: %v)", pdfAddress, parseError)
		// Count it as skipped because nothing was requested from the server.
		downloader.statistics.pdfsSkipped++
		// Stop here; no network traffic happens for non-PDF addresses.
		return
	}
	// Work out the name the server-side file has.
	originalFilename := deriveOriginalFilename(pdfAddress)
	// Convert it to the safe naming scheme.
	sanitizedFilename := sanitizeFilename(originalFilename)
	// Make sure two different PDFs in this run never share a name.
	finalFilename := downloader.reserveUniqueFilename(sanitizedFilename)
	// Log the name conversion.
	writeDebugLog("Original filename: %q | Sanitized filename: %q | Final filename: %q", originalFilename, sanitizedFilename, finalFilename)
	// Build the full path inside the PDFs directory (the only place PDFs are saved).
	finalFilePath := filepath.Join(downloader.outputDirectoryAbsolutePath, finalFilename)
	// Check whether the target already exists from an earlier run.
	if _, statError := os.Stat(finalFilePath); statError == nil {
		// Never overwrite: skip and count it.
		writeDebugLog("Skipping %s because %s already exists", pdfAddress, finalFilePath)
		// Update the statistics.
		downloader.statistics.pdfsSkipped++
		// Nothing more to do.
		return
	} else if !os.IsNotExist(statError) {
		// An unexpected file-system error means we cannot safely write this file.
		writeDebugLog("Cannot check %s: %v. Counting %s as failed", finalFilePath, statError, pdfAddress)
		// Update the statistics.
		downloader.statistics.pdfsFailed++
		// Nothing more to do.
		return
	}
	// Temporary download location (still inside the PDFs directory).
	temporaryFilePath := finalFilePath + ".part"
	// Download with retry logic around a single attempt.
	downloadError := downloader.retryOperation("download "+pdfAddress, func() error {
		// One attempt.
		return downloader.downloadOnce(pdfAddress, temporaryFilePath, finalFilePath)
	})
	// Handle the outcome.
	if downloadError != nil {
		// Report the failure and count it.
		writeDebugLog("Download FAILED: %s -> %v", pdfAddress, downloadError)
		// Update the statistics.
		downloader.statistics.pdfsFailed++
		// Nothing more to do.
		return
	}
	// Report success and count it.
	writeDebugLog("Download completed: %s -> %s", pdfAddress, finalFilePath)
	// Update the statistics.
	downloader.statistics.pdfsDownloaded++
}

// ---------------------------------------------------------------------------------------
// Page crawling
// ---------------------------------------------------------------------------------------

// visitPage fetches one page and returns its parsed URL and all candidate links.
// The boolean result is false when the page could not be fetched.
func (downloader *pdfDownloader) visitPage(pageAddress string) (*url.URL, []string, bool) {
	// Parse the address for later link resolution.
	pageURL, parseError := url.Parse(pageAddress)
	// An unparsable address cannot be visited.
	if parseError != nil {
		// Log, count and stop.
		writeDebugLog("Cannot parse page address %q: %v", pageAddress, parseError)
		// Update the statistics.
		downloader.statistics.pageFetchFailures++
		// Signal failure.
		return nil, nil, false
	}
	// Download the HTML.
	pageBody, fetchError := downloader.fetchPageBody(pageAddress)
	// Failed pages are logged and counted but do not stop the program.
	if fetchError != nil {
		// Log the error.
		writeDebugLog("Page fetch FAILED: %s -> %v", pageAddress, fetchError)
		// Update the statistics.
		downloader.statistics.pageFetchFailures++
		// Signal failure.
		return nil, nil, false
	}
	// Extract every link from the HTML.
	allLinks := extractCandidateLinks(pageBody, pageURL)
	// Log the amount of data found.
	writeDebugLog("Page %s: %d bytes of HTML, %d unique links", pageAddress, len(pageBody), len(allLinks))
	// Return the results.
	return pageURL, allLinks, true
}

// registerDiscoveredPdfs records newly discovered PDF URLs and appends the new ones to the
// category's download queue. PDFs seen earlier (on any page) are not queued twice.
func (downloader *pdfDownloader) registerDiscoveredPdfs(pdfLinks []string, downloadQueue []string) []string {
	// Look at every PDF link of the page.
	for _, pdfLink := range pdfLinks {
		// Skip PDFs that were already discovered earlier.
		if downloader.discoveredPdfUrls[pdfLink] {
			// Log the duplicate.
			writeDebugLog("PDF already discovered, not queuing again: %s", pdfLink)
			// Move on.
			continue
		}
		// Remember the PDF globally.
		downloader.discoveredPdfUrls[pdfLink] = true
		// Count it as a unique discovery.
		downloader.statistics.uniquePdfUrlsDiscovered++
		// Queue it for downloading.
		downloadQueue = append(downloadQueue, pdfLink)
	}
	// Return the extended queue.
	return downloadQueue
}

// processCategory visits one category URL (following pagination), falls back to product
// pages when the category page has no direct PDF links, and downloads all new PDFs.
func (downloader *pdfDownloader) processCategory(categoryAddress string) {
	// Log which category is being processed.
	writeDebugLog("=== Visiting category page: %s ===", categoryAddress)
	// Normalize the starting address.
	startURL, parseError := url.Parse(categoryAddress)
	// Report an invalid category address.
	if parseError != nil {
		// Log, count and stop.
		writeDebugLog("Invalid category address %q: %v", categoryAddress, parseError)
		// Update the statistics.
		downloader.statistics.pageFetchFailures++
		// Stop this category.
		return
	}
	// Normalize (lowercase host, remove fragment).
	normalizedStart, normalizeError := normalizeURL(categoryAddress, startURL)
	// Report a normalization failure.
	if normalizeError != nil {
		// Log, count and stop.
		writeDebugLog("Cannot normalize category address: %v", normalizeError)
		// Update the statistics.
		downloader.statistics.pageFetchFailures++
		// Stop this category.
		return
	}
	// Queue of listing pages still to visit (grows when pagination links are found).
	pagesToVisit := []string{normalizedStart}
	// PDFs found for this category that still need downloading.
	var pdfDownloadQueue []string
	// Product pages found on the listing pages.
	var productPageAddresses []string
	// Set used to avoid listing the same product page twice.
	productPageAlreadyListed := make(map[string]bool)
	// Count of listing pages visited, used for the safety limit.
	listingPagesVisited := 0
	// Process listing pages until none remain or the safety limit is reached.
	for len(pagesToVisit) > 0 && listingPagesVisited < maximumPaginationPagesPerCategory {
		// Take the first page from the queue.
		currentAddress := pagesToVisit[0]
		// Remove it from the queue.
		pagesToVisit = pagesToVisit[1:]
		// Skip pages already visited (prevents pagination loops).
		if downloader.visitedPageUrls[currentAddress] {
			// Log the skip.
			writeDebugLog("Page already visited, skipping: %s", currentAddress)
			// Move on.
			continue
		}
		// Remember the page as visited.
		downloader.visitedPageUrls[currentAddress] = true
		// Count the attempt towards the safety limit.
		listingPagesVisited++
		// Fetch the page and extract its links.
		currentURL, allLinks, pageFetched := downloader.visitPage(currentAddress)
		// A failed page has already been logged and counted.
		if !pageFetched {
			// Move on to the next queued page.
			continue
		}
		// Count the successful category page.
		downloader.statistics.categoryPagesProcessed++
		// Separate the links by purpose.
		pdfLinks := filterPdfLinks(allLinks)
		// Log how many PDF links this page holds.
		writeDebugLog("Category page %s: %d PDF link(s) discovered", currentAddress, len(pdfLinks))
		// Queue new PDFs.
		pdfDownloadQueue = downloader.registerDiscoveredPdfs(pdfLinks, pdfDownloadQueue)
		// Queue further listing pages (pagination).
		for _, paginationLink := range filterPaginationLinks(allLinks, currentURL) {
			// Only add pages that have not been visited yet.
			if !downloader.visitedPageUrls[paginationLink] {
				// Add to the visit queue.
				pagesToVisit = append(pagesToVisit, paginationLink)
			}
		}
		// Remember product detail pages for the fallback step.
		for _, productLink := range filterProductLinks(allLinks, currentURL) {
			// Avoid duplicates.
			if !productPageAlreadyListed[productLink] {
				// Mark as listed.
				productPageAlreadyListed[productLink] = true
				// Remember it.
				productPageAddresses = append(productPageAddresses, productLink)
			}
		}
	}
	// Fallback: visit product pages only when the listing pages exposed no PDFs directly.
	if len(pdfDownloadQueue) == 0 && len(productPageAddresses) > 0 {
		// Log the decision.
		writeDebugLog("No direct PDFs on category pages; inspecting up to %d product page(s)", maximumProductPagesPerCategory)
		// Visit each product page, up to the safety limit.
		for productIndex, productAddress := range productPageAddresses {
			// Stop at the safety limit.
			if productIndex >= maximumProductPagesPerCategory {
				// Log that remaining pages are ignored.
				writeDebugLog("Product page limit reached; ignoring %d further product page(s)", len(productPageAddresses)-productIndex)
				// Leave the loop.
				break
			}
			// Skip product pages already visited through another category.
			if downloader.visitedPageUrls[productAddress] {
				// Log the skip.
				writeDebugLog("Product page already visited, skipping: %s", productAddress)
				// Move on.
				continue
			}
			// Remember the page as visited.
			downloader.visitedPageUrls[productAddress] = true
			// Fetch the product page and extract its links.
			_, productLinks, productFetched := downloader.visitPage(productAddress)
			// A failed page has already been logged and counted.
			if !productFetched {
				// Move on.
				continue
			}
			// Count the successful product page.
			downloader.statistics.productPagesVisited++
			// Keep only PDF links.
			productPdfLinks := filterPdfLinks(productLinks)
			// Log how many PDF links this product page holds.
			writeDebugLog("Product page %s: %d PDF link(s) discovered", productAddress, len(productPdfLinks))
			// Queue new PDFs.
			pdfDownloadQueue = downloader.registerDiscoveredPdfs(productPdfLinks, pdfDownloadQueue)
		}
	}
	// Log how many downloads are about to start for this category.
	writeDebugLog("Category %s: %d new PDF(s) to download", categoryAddress, len(pdfDownloadQueue))
	// Download every queued PDF.
	for _, pdfAddress := range pdfDownloadQueue {
		// Download (or skip) one PDF.
		downloader.downloadPdf(pdfAddress)
	}
}

// ---------------------------------------------------------------------------------------
// Program entry
// ---------------------------------------------------------------------------------------

// printSummary writes the final statistics to the console.
func (downloader *pdfDownloader) printSummary() {
	// Short alias for readability inside this function.
	statistics := downloader.statistics
	// Print each requested figure as its own log line.
	writeDebugLog("================ SUMMARY ================")
	// Category pages.
	writeDebugLog("Category pages processed:      %d", statistics.categoryPagesProcessed)
	// Product pages.
	writeDebugLog("Product pages visited:         %d", statistics.productPagesVisited)
	// Pages that failed.
	writeDebugLog("Pages that could not be read:  %d", statistics.pageFetchFailures)
	// Unique PDF URLs.
	writeDebugLog("Unique PDF URLs discovered:    %d", statistics.uniquePdfUrlsDiscovered)
	// Successful downloads.
	writeDebugLog("PDFs downloaded successfully:  %d", statistics.pdfsDownloaded)
	// Skipped PDFs.
	writeDebugLog("PDFs skipped (already exist):  %d", statistics.pdfsSkipped)
	// Failed downloads.
	writeDebugLog("PDFs failed:                   %d", statistics.pdfsFailed)
	// Output directory.
	writeDebugLog("PDFs directory (absolute path): %s", downloader.outputDirectoryAbsolutePath)
}

// run contains the whole program flow and returns the process exit code.
// Keeping this separate from main lets deferred cleanup run before os.Exit.
func run() int {
	// Log timestamps with microseconds to make request timing visible.
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	// Log the startup.
	writeDebugLog("Application starting")
	// The directory must exist before any download is attempted.
	outputDirectoryAbsolutePath, directoryError := ensureOutputDirectory(outputDirectoryName)
	// Without an output directory nothing can be saved, so stop.
	if directoryError != nil {
		// Report the fatal error.
		writeDebugLog("Fatal error: %v", directoryError)
		// Non-zero exit code signals failure.
		return 1
	}
	// Create the downloader with its shared state.
	downloader := &pdfDownloader{
		httpClient:                  &http.Client{Timeout: httpRequestTimeout}, // Timeout protects against hanging connections.
		outputDirectoryAbsolutePath: outputDirectoryAbsolutePath,               // Where PDFs are saved.
		discoveredPdfUrls:           make(map[string]bool),                     // Starts empty.
		visitedPageUrls:             make(map[string]bool),                     // Starts empty.
		filenamesClaimedThisRun:     make(map[string]bool),                     // Starts empty.
	}
	// Visit every category page in order.
	for _, categoryAddress := range categoryPageAddresses {
		// Process one category.
		downloader.processCategory(categoryAddress)
	}
	// Print the final numbers.
	downloader.printSummary()
	// Log the shutdown.
	writeDebugLog("Application finished")
	// Report failure through the exit code when anything went wrong.
	if downloader.statistics.pdfsFailed > 0 || downloader.statistics.pageFetchFailures > 0 {
		// Non-zero exit code signals partial failure.
		return 1
	}
	// Zero means everything succeeded.
	return 0
}

// main is the program entry point.
func main() {
	// Run the program and exit with its status code.
	os.Exit(run())
}
