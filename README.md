# 📚 detecto-com-documentation

> 🧰 A tidy, offline-ready archive of **Detecto** product documentation (manuals, bulletins, certifications, cleaning guides and more), collected with a polite, rate-limited downloader. ⚖️🏥🍎🐾🚚

![Go](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go&logoColor=white)
![Python](https://img.shields.io/badge/Python-3.x-3776AB?logo=python&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)
![Docs](https://img.shields.io/badge/Content-PDF%20Documentation-red)
![Status](https://img.shields.io/badge/Status-Active-brightgreen)

---

## 📖 Table of Contents

1. [🌟 What Is This?](#-what-is-this)
2. [🧭 At a Glance (Key Facts)](#-at-a-glance-key-facts)
3. [🏢 About Detecto and Its Product World](#-about-detecto-and-its-product-world)
4. [📦 What's Inside the Archive?](#-whats-inside-the-archive)
5. [🔬 Sample Data: What a Real Product Page Contains](#-sample-data-what-a-real-product-page-contains)
6. [🗃️ Data Dictionary (Fields You Will Meet)](#️-data-dictionary-fields-you-will-meet)
7. [💎 Why This Data Is Valuable](#-why-this-data-is-valuable)
8. [🧑‍💼 Use-Case Playbooks by Role](#-use-case-playbooks-by-role)
9. [📚 Glossary of Terms](#-glossary-of-terms)
10. [🗂️ Repository Layout](#️-repository-layout)
11. [🚀 Quick Start](#-quick-start)
12. [⚙️ How the Downloader Works](#️-how-the-downloader-works)
13. [🔧 Configuration Reference](#-configuration-reference)
14. [🔤 Filename Rules](#-filename-rules)
15. [🧪 Reliability Features](#-reliability-features)
16. [🖥️ Example Output](#️-example-output)
17. [🛠️ Troubleshooting](#️-troubleshooting)
18. [❓ FAQ](#-faq)
19. [🔎 Searching and Analyzing the Archive](#-searching-and-analyzing-the-archive)
20. [🔐 Ethics, Safety and Legal Notes](#-ethics-safety-and-legal-notes)
21. [🗺️ Roadmap Ideas](#️-roadmap-ideas)
22. [🤝 Contributing](#-contributing)
23. [📄 License](#-license)

---

## 🌟 What Is This?

**Detecto** is a well-known maker of precision weighing equipment: 🍼 baby scales, 🪑 chair scales, 🩺 physician scales, ♿ wheelchair scales, 🐕 veterinary scales, 🍞 food-service scales, 📦 shipping scales and more. 🏥

Every product page on the Detecto website publishes **PDF documents**: owner's manuals, bulletins, declarations of conformity, cleaning guides and similar material. They are spread across many pages, sometimes available in several languages, and easy to lose track of. 😅

This repository solves that by:

1. 🕷️ **Visiting** selected Detecto product pages.
2. 🔎 **Finding** every PDF link, including links hidden inside the page's embedded data.
3. ⬇️ **Downloading** each PDF safely, with validation.
4. 🧹 **Renaming** files to clean, predictable, lowercase names.
5. 📁 **Storing** everything together in the `PDFs/` folder.

The result is one searchable, organized, offline-friendly **documentation library**. 🎉

---

## 🧭 At a Glance (Key Facts)

| 🏷️ Item                 | 📋 Details                                                               |
| ----------------------- | ------------------------------------------------------------------------ |
| **Purpose**             | Build a local library of Detecto product documentation                   |
| **Source**              | Public Detecto product pages (see the target list below)                 |
| **Output folder**       | `PDFs/`                                                                  |
| **Languages**           | Go (standard library only) and Python                                    |
| **Document formats**    | PDF                                                                      |
| **Filename style**      | `lowercase_with_single_underscores.pdf`                                  |
| **Validation**          | Non-empty file + `%PDF-` signature + `.pdf` address + content-type check |
| **Politeness**          | Request spacing, timeouts, retries with back-off                         |
| **License (code)**      | MIT                                                                      |
| **License (documents)** | Property of the original rights holders                                  |

### 🎯 Target Pages

The downloader is set up for these six baby-scale product pages on `detecto.com`:

|  #  | Product page (slug)     | Product                          |
| :-: | ----------------------- | -------------------------------- |
|  1  | `6745-digital`          | 6745 Digital Baby Scale          |
|  2  | `mb130-digital`         | MB130 Digital Baby Scale         |
|  3  | `8440-digital`          | 8440 Digital Baby Scale          |
|  4  | `8435-digital`          | 8435 Digital Baby Scale          |
|  5  | `450-series-weigh-beam` | 450 Series Weigh Beam Baby Scale |
|  6  | `hs25kgp-suspended`     | HS25KGP Suspended Baby Scale     |

Full pattern: `https://detecto.com/category/baby-scales/subcategory/<slug>`

> 💡 Although the address says "category", each of these pages is really a **single product page** with tabs for Overview, Models, Features, Specifications and Resources. The PDFs live in the **Resources** data.

---

## 🏢 About Detecto and Its Product World

Understanding how the company organizes its catalog makes the archive much easier to navigate. 🗺️

### 🧱 Eleven Product Groups

| 🧩 Product group                             | 🗂️ Categories inside                                                                               |
| -------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| 🍼 **Baby Scales**                           | Baby Scales                                                                                        |
| 🪑 **Chair Scales**                          | Chair Scales                                                                                       |
| 🛏️ **Patient Lift / In-Bed Scales**          | In-Bed Scales, Patient Lift Scales                                                                 |
| 🧍 **Stand-On Scales**                       | Bariatric, Physician, Athletic, Home Healthcare, Stadiometers                                      |
| 🧷 **Wet Diaper / Chux Scales**              | Chux Underpad Scales, Wet Diaper Scales                                                            |
| ♿ **Wheelchair / Floor / Stretcher Scales** | Wheelchair Scales, In-Floor Dialysis Scales, Stretcher Scales                                      |
| 🛒 **Medical Carts**                         | Medical Carts                                                                                      |
| 🗑️ **Waste Receptacles / Glove Box Holders** | Glove Box Holders, Waste Receptacles                                                               |
| 🍽️ **Food Service Scales**                   | Baker Dough, Hanging, Ingredient, Portion, Price Computing, Retail POS/Logistics, Top Loading Dial |
| 🔬 **Precision Balance Scales**              | Precision Digital Balances                                                                         |
| 🐾 **Veterinary Scales**                     | Veterinary Scales                                                                                  |

### 🏭 Twenty-Two Industries Served

🚑 Acute Care · 👵 Assisted Living & Senior Care · 🏃 Athletics · 🥖 Bakeries · ⚖️ Bariatrics · 🍲 Cafeterias · 🏥 Clinics · 👩‍🍳 Commercial Kitchens · 💧 Dialysis · 🎣 Fishing · 🛒 Grocery Stores · 🏠 Home Healthcare · 🏨 Long-Term Care · 🧑‍🏭 Occupational Health · 👶 Pediatrics · 🩺 Primary Care · 🦿 Rehabilitation & Physical Therapy · 🧾 Retail Point-of-Sale · 📦 Shipping/Logistics · 🔪 Surgery Centers · ⏱️ Urgent Care · 🐶 Veterinary

> 📌 The wide range shows why documentation matters: the same brand of equipment ends up in NICUs, bakeries, fishing boats and warehouses, each with different cleaning, accuracy and compliance needs.

### 🗃️ Resource Types Published by the Site

The website groups downloadable material into types such as:

- 📘 **Manuals**
- 📰 **Brochures** (product bulletins)
- 🏅 **Certifications** (for example, declarations of conformity)
- 🎬 **Videos**
- 🖼️ **Logos**
- 💾 **Software** and 🧊 **3D model** fields exist as well

This archive focuses on the **PDF** items inside those resources.

---

## 📦 What's Inside the Archive?

Typical document types you will find in `PDFs/`:

| Emoji | Document type                        | What it usually contains                                             |
| :---: | ------------------------------------ | -------------------------------------------------------------------- |
|  📘   | **Owner's manuals**                  | Unpacking, setup, operation, care and cleaning, parts lists          |
|  📰   | **Product bulletins / brochures**    | Features, capacities, options, model comparisons                     |
|  📋   | **Specification sheets**             | Capacity and resolution, display, power, dimensions, shipping weight |
|  🏅   | **Declarations of conformity**       | CE and regulatory compliance statements                              |
|  🧼   | **Cleaning and disinfecting guides** | Approved cleaning methods for scales and indicators                  |
|  🌍   | **Translated documents**             | For example, Spanish editions of manuals and bulletins               |
|  💳   | **Company documents**                | For example, a Letter of Credit PDF linked from the site footer      |

> 💡 The exact file list depends on which pages were downloaded and when. Browse `PDFs/` to see what is currently in the archive.

---

## 🔬 Sample Data: What a Real Product Page Contains

To show the kind of information this archive captures, here are real examples taken from public Detecto product pages. 📊

### 🪝 Example 1: HS25KGP Suspended Baby Scale

A lightweight, portable **mechanical** hanging scale for remote and field healthcare. 🏕️

| 📏 Specification    | 📋 Value                                             |
| ------------------- | ---------------------------------------------------- |
| Capacity × division | **25 kg × 100 g**                                    |
| Display             | 8 in (20 cm) shatterproof glass dial                 |
| Power               | None required (no batteries, no AC) 🔋🚫             |
| Weighing seat       | Waterproof vinyl weighing trousers (suspended sling) |
| Hook                | 4 in (10 cm) zinc-plated S-hook                      |
| Shipping weight     | 4 lb (1.8 kg)                                        |
| Best for            | Remote clinics, field work, mobile care              |

**📄 PDFs published on this product's page (as of the page data I reviewed):**

| 📘 Document                                               | 🏷️ Type       | 📦 Approx. size | 🌐 Language |
| --------------------------------------------------------- | ------------- | --------------: | :---------: |
| HS25KGP Owner's Manual                                    | Manual        |         ~103 KB |     EN      |
| HS25KGP Manual (Spanish)                                  | Manual        |         ~140 KB |     ES      |
| HS25KGP Pediatric Scale Bulletin (English)                | Brochure      |         ~454 KB |     EN      |
| HS25KGP Pediatric Scale Bulletin (Spanish)                | Brochure      |         ~464 KB |     ES      |
| HS25KGP Declaration of Conformity                         | Certification |         ~266 KB |     EN      |
| Cleaning and Disinfecting Your Scale and Weight Indicator | Manual        |         ~557 KB |     EN      |

> 🧠 **Observation:** the cleaning and disinfecting guide is shared by **many** products. The page data shows it attached to dozens of product entries (baby, chair, physician, wheelchair, stretcher, veterinary scales and more). That is exactly why **de-duplication** matters: one PDF, many pages. ♻️

### 🍼 Example 2: 6745 Digital Baby Scale

A high-precision **digital** infant scale for NICUs, nurseries, delivery rooms and pediatric practices. 👶

| 📏 Specification           | 📋 Value                                               |
| -------------------------- | ------------------------------------------------------ |
| Capacity                   | **30 lb (15 kg)**                                      |
| Resolution                 | **1 gram** divisions                                   |
| Display                    | Full-color 3.2-inch touchscreen 🖥️                     |
| Length measurement         | Integrated measuring tape up to 22 in (56 cm) 📏       |
| Patient ID                 | 8-digit patient ID entry for EMR/EHR workflows         |
| Tracking                   | Built-in Milk and Diaper tracking keys 🍼              |
| Memory                     | Stores up to five patient weights                      |
| Connectivity               | Wi-Fi, Bluetooth® 5.0, RS-232, USB-B 📡                |
| Power                      | Six C-cell batteries **or** the included AC adapter 🔌 |
| Certification / compliance | NTEP-certified, TAA-compliant                          |
| Manufacturing              | Made in the USA 🇺🇸                                     |
| Unit lock                  | OneWeigh® technology locks the display to lb or kg     |

### 🔍 What These Examples Teach Us

- 📈 Specifications vary hugely, from a battery-free dial scale to a connected touchscreen scale.
- 🧾 Documents come in **families**: manual + bulletin + conformity statement + shared cleaning guide.
- 🌍 Some documents exist in **multiple languages**.
- 🔗 Many documents are **shared across products**, so a good downloader must avoid duplicates.

---

## 🗃️ Data Dictionary (Fields You Will Meet)

Detecto's website delivers product data in a structured format behind the scenes. If you inspect the page data, these are the useful fields. 🧬

### 🧩 Product ("subcategory") record

| Field                              | Meaning                                                   |
| ---------------------------------- | --------------------------------------------------------- |
| `title`                            | Product display name                                      |
| `slug`                             | URL-friendly identifier (for example `hs25kgp-suspended`) |
| `overview`                         | Long descriptive text                                     |
| `features`                         | Feature bullet lists                                      |
| `specifications`                   | Specification bullet lists                                |
| `galleryImages` / `hero` / `image` | Product pictures                                          |
| `categories`                       | Which category and group the product belongs to           |
| `resources`                        | List of attached documents, videos and images             |
| `discontinued`                     | Whether the product is discontinued                       |
| `typeCounts`                       | Number of models, options and parts                       |

### 📎 Resource record

| Field                | Meaning                                                |
| -------------------- | ------------------------------------------------------ |
| `title`              | Human-readable document title                          |
| `type`               | Manuals, Brochures, Certifications, Videos, Logos, ... |
| `sub_category`       | The products this resource is attached to              |
| `pdfs[]`             | One or more attached PDF entries                       |
| `pdfs[].description` | Short explanation of the document                      |
| `pdfs[].file.name`   | Original filename, such as `HS25KGP_Owners.pdf`        |
| `pdfs[].file.url`    | Direct storage address of the PDF                      |
| `pdfs[].file.size`   | File size (in kilobytes)                               |
| `pdfs[].file.mime`   | MIME type, `application/pdf` for PDFs                  |

### 🔗 Where the PDF files are stored

PDF files are served from a separate **object-storage host**, not from `detecto.com` itself, under a folder like:

```text
.../prod/Resource/PDFs/<Name>_<short-hash>.pdf
```

The short hash at the end (for example `_5f3a162c95`) is added by the storage system. 🔑 It is why saved filenames can end with a hash, and it is also helpful: the hash changes when a file changes.

---

## 💎 Why This Data Is Valuable

This is the heart of the project. 🫀 Here is why a well-organized documentation archive is worth having.

### 📘 1. Manuals

Manuals are the **most useful** document type for anyone who owns or maintains a scale.

- 🧭 **Step-by-step operation**: calibration, taring, unit switching, memory functions.
- 🔧 **Setup instructions**: assembly, leveling, power connection, battery installation.
- 🚨 **Troubleshooting**: error codes and what they mean.
- 📦 **Unpacking guidance**: how to inspect and install equipment correctly.
- 🧾 **Parts lists**: exact part names and numbers for ordering replacements.
- 🔋 **Power guidance**: battery types, adapters and charging notes.
- 🧽 **Care instructions**: what to clean with, and what to avoid.

> 🕰️ Manufacturers reorganize websites and retire pages. Having the manual **saved locally** means you are never stuck when a link breaks.

### 🛠️ 2. Repair and Maintenance

- 🧩 **Identify the right part** using the parts lists and diagrams.
- 🔩 **Follow factory-approved procedures** instead of guessing.
- 📅 **Plan preventive maintenance**: cleaning intervals, inspection points, calibration checks.
- 💸 **Save money**: diagnose simple issues (loose connection, dead battery, wrong setting) before calling a technician.
- ⏱️ **Reduce downtime**: especially important in hospitals, clinics and kitchens. 🏥🍳
- 🧾 **Build a maintenance log** that references the exact manual revision you used.

### 🔌 3. Integration and Interoperability

Many modern scales connect to other systems: 🖥️ EMR/EHR software, 🧾 point-of-sale systems, 🖨️ printers, 💻 PCs.

Documentation often describes:

- 🔗 **Connectivity options**: USB, RS-232, Wi-Fi and Bluetooth.
- 📡 **Output formats and settings** (when published).
- ⚡ **Power and electrical requirements**.
- 🧱 **Mounting and installation constraints**.
- 🏷️ **Patient or item ID entry** capabilities.

This helps developers, integrators and IT teams build reliable connections between scales and software. 👩‍💻👨‍💻

### 🔍 4. Reverse Engineering (the Legitimate Kind)

"Reverse engineering" has a bad reputation, but **analyzing documentation to understand how a product behaves** is a normal, legal and valuable engineering activity when done for the right reasons. ✅

Good uses include:

- 🧠 **Understanding product architecture** from specifications, option lists and part breakdowns.
- 🔄 **Interoperability work**: writing software or hardware that works _with_ a scale rather than copying it.
- 🧪 **Test planning**: building QA checklists from documented features and limits.
- 📊 **Competitive analysis**: comparing documented capacities, resolutions, certifications and features across product lines.
- 🏗️ **Design inspiration**: learning how professional-grade equipment is specified (not cloning it).
- 🧯 **Failure analysis**: matching a malfunction to documented operating limits.
- 🧬 **Product-line mapping**: seeing how models share parts, options and accessories.

> ⚠️ **Please remember:** documentation analysis is for learning, integration and maintenance. Respect patents, trademarks, copyrights and the manufacturer's terms. Do not use this material to counterfeit products or misrepresent them. 🙏

### 📐 5. Blueprints, Diagrams and Specifications

Be realistic: these PDFs are **public product literature**, not confidential factory CAD files. 🙂 Even so, they often include very practical technical content:

- 📏 **Dimensions and footprints**, so you can plan counter space, carts and wall mounts.
- 🧱 **Platform sizes and capacities** for choosing the right model.
- 🔍 **Exploded views and parts diagrams** (where included).
- 🧮 **Capacity and resolution tables** (for example, 25 kg × 100 g).
- ⚖️ **Weight and shipping data** for logistics planning. 🚚
- 🪛 **Installation notes** for mounting, leveling and wiring.

Think of them as **practical blueprints for deployment**: they tell you what will fit, what it needs and how it is built to be used. 🗺️

### 🏅 6. Compliance and Certification

- 📜 **Declarations of conformity** document regulatory compliance (for example, CE).
- ✅ **Audit readiness**: quickly produce paperwork when inspectors or quality teams ask.
- 🏥 **Procurement checks**: confirm a scale meets your facility's requirements before buying.
- 🧾 **Record keeping**: keep a dated copy of the documents that applied when you bought the device.
- 🏷️ **Spec-sheet claims** such as NTEP certification and TAA compliance can be verified against the documents.

### 🧼 7. Cleaning and Infection Control

- 🧴 **Approved cleaning and disinfecting methods** protect both the equipment and your patients.
- 🏥 **Healthcare-friendly**: infection prevention teams can reference the exact guidance.
- 🚫 **Avoid damage**: using the wrong chemicals can harm displays, seals and finishes.
- 📋 **Standardize protocols** across wards and facilities using one trusted source.

### 🎓 8. Training and Education

- 🧑‍🏫 **Onboarding new staff**: nurses, vet techs, kitchen staff and warehouse teams.
- 📚 **Course material**: biomedical engineering, metrology and instrumentation classes.
- 🖨️ **Quick-reference sheets**: print the pages people actually need.
- 🌐 **Multilingual teams**: translated manuals reduce mistakes. 🇪🇸

### 🛒 9. Procurement and Product Comparison

- 🆚 **Compare models side by side**: capacity, resolution, power, connectivity and features.
- 💰 **Make informed buying decisions** before requesting quotes.
- 🧮 **Check accessories and options** so nothing is missed at purchase time.
- 📈 **Standardize equipment** across multiple clinics, stores or sites. 🏢

### 🤖 10. Data Science, Search and AI

A clean document set is a great foundation for modern tooling:

- 🔎 **Full-text search** across every manual at once.
- 🧠 **Retrieval-augmented assistants** (RAG) that answer questions such as "How do I clean this scale?" from the real manual.
- 🏷️ **Information extraction**: pull capacities, dimensions and model numbers into a spreadsheet.
- 🗂️ **Classification experiments** on technical documents.
- 🌍 **Translation and accessibility** projects (for example, large-print versions where permitted).
- 📉 **Change tracking**: compare new downloads with old ones to see what the manufacturer updated.

### 🗄️ 11. Preservation and Offline Access

- 📴 **Works with no internet**: ideal for remote clinics, field work and shipboard use. ⛵
- 🕰️ **Protects against link rot** and removed pages.
- 🧯 **Disaster recovery**: documentation is available even if a website is down.
- 📁 **Version your own copies** with Git so you can see when files were added.

### 🧾 12. Warranty, Support and Dispute Resolution

- 📞 **Describe your model precisely** when contacting support.
- 🧷 **Prove what the documentation said** at the time of purchase or installation.
- 🔁 **Return and warranty conversations** go smoother with the exact manual on hand.

### 🌍 13. Accessibility and Field Work

- 📱 PDFs open on almost any phone or tablet.
- 🏕️ Offline copies matter where connectivity is poor, such as the remote and field settings the HS25KGP is designed for.
- 🔤 Text-based PDFs can be read aloud by assistive tools, copied into translators or reformatted.

---

## 🧑‍💼 Use-Case Playbooks by Role

Quick, practical ideas for different people. 🎭

### 🩺 Biomedical / Clinical Engineer

1. 📥 Download the manual and conformity statement for each scale model in your hospital.
2. 🧾 Build an equipment sheet: model, capacity, resolution, power type, connectivity.
3. 🛠️ Schedule maintenance using the manual's care instructions.
4. 🧼 Share the cleaning guide with infection prevention.

### 👩‍💻 Software Developer / Integrator

1. 🔌 Search the documents for RS-232, USB, Bluetooth or Wi-Fi sections.
2. 🧪 Build a test plan around documented capacities, units and limits.
3. 🧱 Confirm electrical and mounting constraints before designing hardware.
4. 🔄 Keep your integration notes linked to the manual revision you used.

### 🛒 Purchasing / Procurement Manager

1. 📋 Extract capacity, resolution and certifications into a comparison table.
2. 🏅 Verify claimed certifications against the documents.
3. 🧮 Plan accessories (carts, adapters, batteries) from the option lists.

### 👩‍🏫 Teacher / Trainer

1. 🖨️ Print only the operation pages for a laminated quick-start card.
2. 🌐 Use Spanish editions for bilingual teams where available.
3. 🎓 Use real manuals as case studies for technical-reading classes.

### 🧑‍🔧 Field Technician

1. 📱 Keep the archive on a tablet or phone for offline use.
2. 🔩 Look up parts lists before arriving at the site.
3. 📸 Photograph the serial label and match it to the correct document.

### 📊 Data Analyst / Researcher

1. 🗂️ Convert PDFs to text and build a searchable index.
2. 🧬 Extract model numbers, capacities and resolutions into a dataset.
3. 📈 Track how documentation changes over time.

### 🐾 Veterinary Practice Manager

1. 🐶 Keep manuals for every animal scale in the clinic.
2. 🧼 Use cleaning guidance between patients.
3. 📅 Track calibration and maintenance dates alongside the documents.

---

## 📚 Glossary of Terms

| 🔤 Term                            | 📖 Meaning                                                                                                        |
| ---------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| **Capacity**                       | The maximum weight a scale can measure                                                                            |
| **Resolution / division**          | The smallest weight step the scale shows (for example, 1 g or 100 g)                                              |
| **Tare**                           | Setting the scale to zero with a container or accessory on it so only the contents are weighed                    |
| **NTEP**                           | U.S. National Type Evaluation Program, a certification used for scales that must meet legal measurement standards |
| **Legal-for-trade**                | A scale approved for use where the weight determines a price or sale                                              |
| **TAA-compliant**                  | Meets the Trade Agreements Act requirements for certain U.S. government purchases                                 |
| **CE / Declaration of Conformity** | A statement that a product meets applicable European requirements                                                 |
| **EMR / EHR**                      | Electronic Medical Record / Electronic Health Record systems                                                      |
| **RS-232**                         | A traditional serial connection used to send data from devices to computers or printers                           |
| **USB-B**                          | A type of USB connector commonly used on devices connecting to a host                                             |
| **Bluetooth® 5.0**                 | A wireless standard for short-range device connections                                                            |
| **Weigh beam**                     | A mechanical balance scale using sliding weights on a beam                                                        |
| **Suspended / hanging scale**      | A scale that hangs from a hook and weighs what is suspended from it                                               |
| **Dial scale**                     | A mechanical scale that shows weight on a dial with a pointer                                                     |
| **Bulletin / brochure**            | A marketing and specification summary of a product                                                                |
| **Parts list**                     | A list of replaceable components with names or numbers                                                            |
| **Link rot**                       | When web links stop working because pages move or disappear                                                       |
| **RAG**                            | Retrieval-augmented generation: an AI approach that answers questions using retrieved documents                   |
| **MIME type**                      | A label telling software what kind of file something is (for example, `application/pdf`)                          |

---

## 🗂️ Repository Layout

```text
detecto-com-documentation/
├── 📁 PDFs/            # Downloaded documentation (clean, lowercase filenames)
├── 🐹 main.go          # Go downloader (standard library only)
├── 🐍 main.py          # Python script
├── 📄 go.mod           # Go module definition
├── 📦 requirements.txt # Python dependencies
├── 🚀 upload.sh        # Helper script for publishing changes
├── 🙈 .gitignore       # Files Git should ignore
├── ⚖️ LICENSE          # MIT license (applies to the code, see notes below)
└── 📘 README.md        # You are here!
```

---

## 🚀 Quick Start

### ✅ Requirements

| Tool          | Version                                               | Needed for           |
| ------------- | ----------------------------------------------------- | -------------------- |
| 🐹 Go         | 1.21 or newer                                         | `main.go`            |
| 🐍 Python     | 3.x                                                   | `main.py`            |
| 🌐 Internet   | any                                                   | Downloading the PDFs |
| 💽 Disk space | a few MB to hundreds of MB depending on how many PDFs | Storing the library  |

### 🐹 Option A: Go

```bash
# 1️⃣ Clone the repository
git clone https://github.com/PrajwalKoirala638/detecto-com-documentation.git
cd detecto-com-documentation

# 2️⃣ Run the downloader (no extra dependencies needed)
go run main.go
```

### 🐍 Option B: Python

```bash
# 1️⃣ (Optional) create a virtual environment
python3 -m venv .venv
source .venv/bin/activate        # Windows: .venv\Scripts\activate

# 2️⃣ Install dependencies
pip install -r requirements.txt

# 3️⃣ Run the script
python main.py
```

After it finishes, open the **`PDFs/`** folder and enjoy your library! 🎊

### 🧭 Handy Commands

```bash
# 📋 List downloaded PDFs
ls -1 PDFs | sort

# 🔢 Count them
ls -1 PDFs/*.pdf | wc -l

# 💽 Total size
du -sh PDFs

# ✅ Quick sanity check: every file should start with %PDF-
for f in PDFs/*.pdf; do head -c 5 "$f" | grep -q '%PDF-' || echo "⚠️ Not a PDF: $f"; done
```

---

## ⚙️ How the Downloader Works

```text
🌐 Category page ──► 🔎 Find links ──► 📄 PDF links? ──► ✅ Validate ──► 💾 Save to PDFs/
                         │                    │
                         │                    └─ no ─► 🧭 Visit product / resource pages
                         └──► 📑 Follow pagination (with safety limits)
```

### 🪜 Step by Step

1. 🧭 Visits each configured page with a browser-style **User-Agent**.
2. 📊 Checks the **HTTP status code** and logs it.
3. 🔎 Extracts links from normal HTML **and** from embedded JSON (modern JavaScript sites often hide links there).
4. 🧹 Normalizes URLs (absolute or relative) and **removes duplicates**.
5. 🧐 Only requests addresses that end in **`.pdf`**.
6. 🛡️ Rejects responses announced as HTML, text or JSON.
7. ✅ Verifies the saved file is **non-empty** and contains the **`%PDF-`** signature.
8. 🏷️ Applies the filename rules below and saves into `PDFs/`.
9. 📋 Prints a **final summary**: pages processed, PDFs discovered, downloaded, skipped and failed.

### 🧠 Why Embedded Data Matters

Modern websites (this one uses a React-based framework) often send page data as JSON inside the HTML, then draw the page in your browser. A basic link scraper that only looks at `<a href>` tags would miss most PDFs. 🙈

This downloader scans:

- 🏷️ attribute values (`href`, `src`, `data-*`)
- 📜 JSON text inside script blocks, including escaped slashes and quotes
- 🔗 absolute and root-relative addresses that end in `.pdf`

### ♻️ De-duplication

The same PDF can appear on many product pages (the cleaning guide is a good example). 🧽 The downloader remembers every normalized PDF address it has seen and downloads each **once**.

### 🧱 Safe Writes

```text
download ──► PDFs/name.pdf.part ──► validate ──► rename to PDFs/name.pdf
                    │                   │
                    └── on any failure: delete the .part file 🧹
```

A half-finished download can never be mistaken for a complete PDF. ✅

---

## 🔧 Configuration Reference

Defaults used by the reference Go downloader (all are constants near the top of `main.go`, so they are easy to change):

| ⚙️ Setting             | 🎛️ Default | 📝 What it controls                     |
| ---------------------- | ---------: | --------------------------------------- |
| Output directory       |     `PDFs` | Where files are saved                   |
| HTTP timeout           |       90 s | Maximum time for one request            |
| Attempts per request   |          4 | How many times a request is tried       |
| First retry delay      |        2 s | Doubles each retry (2 s, 4 s, 8 s)      |
| Maximum retry delay    |       30 s | Upper limit for waiting between retries |
| Delay between requests |     750 ms | Keeps the request rate polite           |
| Max pagination pages   |         20 | Safety limit per category               |
| Max product pages      |         25 | Safety limit per category fallback      |
| Max HTML page size     |     10 MiB | Larger pages are truncated              |
| PDF signature window   | 1024 bytes | Where `%PDF-` must appear               |

### 🪪 User-Agent

Requests use a browser-style User-Agent string:

```text
Mozilla/5.0 (X11; CrOS x86_64 14541.0.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36
```

### ➕ Adding More Pages

To archive more products, add their page addresses to the list of target URLs near the top of the program, then run it again. ✨ Existing files are **not overwritten**, so reruns only add new documents.

---

## 🔤 Filename Rules

All filenames are made **safe and predictable**:

- 🔡 lowercase letters `a-z`, digits `0-9` and underscores `_` only
- 🚫 no leading or trailing underscore
- 🚫 no double underscores
- 🧩 duplicates get a numeric suffix, such as `_2`, so files are never overwritten
- 🆘 if nothing usable is left, a fallback name (`document.pdf`) is used

| Original name                  | Saved as                       |
| ------------------------------ | ------------------------------ |
| `Baby Scale 6745.pdf`          | `baby_scale_6745.pdf`          |
| `MB-130 DIGITAL.pdf`           | `mb_130_digital.pdf`           |
| `_Digital Scale_.pdf`          | `digital_scale.pdf`            |
| `8435__Series!!.pdf`           | `8435_series.pdf`              |
| `HS_25_KGP_DOC_5f3a162c95.pdf` | `hs_25_kgp_doc_5f3a162c95.pdf` |

---

## 🧪 Reliability Features

- ⏱️ **Timeouts** so a slow server cannot hang the program
- 🔁 **Automatic retries** with increasing delays for temporary errors (HTTP 408, 429, 5xx and network problems)
- 🐢 **Polite request spacing** to avoid overloading the website
- 🧱 **Safe writes**: downloads go to a temporary `.part` file and are renamed only after validation
- 🚫 **Never overwrites** existing PDFs
- 📏 **Length check**: if the server announces a size, the received bytes must match
- 🧾 **Detailed debug logging** of requests, status codes, filename changes, retries and errors
- 🔒 **No silent failures**: every error is logged and counted
- 🚦 **Exit code**: non-zero when any page or download fails, handy for scripts and automation

---

## 🖥️ Example Output

A short illustration of what the log looks like (times and counts will differ):

```text
[DEBUG] Application starting
[DEBUG] Created output directory: /home/you/detecto-com-documentation/PDFs
[DEBUG] === Visiting category page: https://detecto.com/category/baby-scales/subcategory/hs25kgp-suspended ===
[DEBUG] HTTP GET https://detecto.com/category/baby-scales/subcategory/hs25kgp-suspended
[DEBUG] HTTP response for ...: 200 OK
[DEBUG] Category page ...: 6 PDF link(s) discovered
[DEBUG] Download starting: https://.../prod/Resource/PDFs/HS_25_KGP_Owners_083fe19d0f.pdf
[DEBUG] Original filename: "HS_25_KGP_Owners_083fe19d0f.pdf" | Sanitized filename: "hs_25_kgp_owners_083fe19d0f.pdf"
[DEBUG] Saved 105861 bytes to .../PDFs/hs_25_kgp_owners_083fe19d0f.pdf
[DEBUG] ================ SUMMARY ================
[DEBUG] Category pages processed:      6
[DEBUG] Product pages visited:         0
[DEBUG] Unique PDF URLs discovered:    ...
[DEBUG] PDFs downloaded successfully:  ...
[DEBUG] PDFs skipped (already exist):  ...
[DEBUG] PDFs failed:                   0
[DEBUG] PDFs directory (absolute path): /home/you/detecto-com-documentation/PDFs
[DEBUG] Application finished
```

---

## 🛠️ Troubleshooting

| 😕 Problem                           | 🔍 Likely cause                                         | ✅ What to try                                                                                                          |
| ------------------------------------ | ------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------- |
| **0 PDFs discovered**                | The page loads its data in a way the scraper cannot see | Check the log line "PDF link(s) discovered"; inspect the page source for `.pdf` addresses; add the correct data address |
| **HTTP 403 / 429**                   | The site is limiting automated traffic                  | Wait, lower the request rate, check the site's terms and `robots.txt`                                                   |
| **"not a PDF" errors**               | A link returned a web page instead of a document        | Open the link in a browser to see what it returns; the check is working as intended                                     |
| **Files named with a trailing hash** | The storage system appends a hash to filenames          | Normal; rename afterwards if you prefer shorter names                                                                   |
| **Skipped counts are high**          | Files already exist from an earlier run                 | Expected; existing files are never overwritten                                                                          |
| **Timeouts**                         | Slow connection or busy server                          | Increase the timeout constant or try again later                                                                        |
| **Permission denied**                | The folder is not writable                              | Run from a folder you own                                                                                               |
| **Certificate or network errors**    | Proxy or firewall in the way                            | Check network settings                                                                                                  |

---

## ❓ FAQ

**❔ Is this an official Detecto repository?**
No. It is an unofficial personal archive and is not affiliated with Detecto. 🙅

**❔ Can I redistribute the PDFs?**
The PDFs belong to their owners. Check the manufacturer's terms before sharing them publicly. ⚖️

**❔ Will the archive stay up to date?**
Only if you rerun the downloader. Manufacturers update documents, so confirm important details on the official site. 🔄

**❔ Why do some filenames end with random characters?**
They come from the storage system's hash. It also helps tell different versions of a file apart. 🔑

**❔ Why only baby scales?**
That is the current target list. You can add more product pages to the list. 🍼➡️🪑

**❔ Do the documents include CAD files or engineering drawings?**
Generally no. They are public product documents, though many include dimensions, parts lists and diagrams. 📐

**❔ Can I use these for repairs?**
They are a helpful reference, but follow the manufacturer's instructions and use qualified technicians for calibration and certified equipment. 🧑‍🔧

**❔ Does the downloader respect `robots.txt`?**
It does not read it automatically. Please review it and the site's terms yourself before running the tool, and keep the request rate low. 🤝

---

## 🔎 Searching and Analyzing the Archive

Some ideas to get more value from the PDFs. 💡

### 🔤 Search for a word across every PDF

```bash
# Install poppler-utils (Linux) or use your platform's equivalent
for f in PDFs/*.pdf; do
  if pdftotext "$f" - 2>/dev/null | grep -qi "calibration"; then
    echo "📄 $f"
  fi
done
```

### 🧾 Convert everything to text

```bash
mkdir -p text
for f in PDFs/*.pdf; do
  pdftotext -layout "$f" "text/$(basename "${f%.pdf}").txt"
done
```

### 🧮 Ideas for a spreadsheet

| Column        | Example                           |
| ------------- | --------------------------------- |
| File name     | `hs_25_kgp_owners_083fe19d0f.pdf` |
| Product       | HS25KGP Suspended Baby Scale      |
| Document type | Manual                            |
| Language      | EN                                |
| Capacity      | 25 kg                             |
| Resolution    | 100 g                             |
| Power         | None                              |
| Notes         | Mechanical dial scale             |

### 🤖 Build a question-answering assistant

1. ✂️ Split each document's text into chunks.
2. 🧬 Create embeddings and store them in a vector database.
3. ❓ Retrieve the best chunks for a question.
4. 💬 Let an AI model answer **only** from the retrieved text and cite the source file.

---

## 🔐 Ethics, Safety and Legal Notes

Please read this section carefully. 🙏

1. ©️ **Documents belong to their owner.** The PDFs are the property of **Detecto / Cardinal Scale** (or their respective rights holders). This repository is an unofficial personal archive and is **not affiliated with or endorsed by Detecto**. 🤝
2. ⚖️ **The MIT license covers the code only** (the downloader scripts), not the downloaded PDFs. Check the manufacturer's terms before redistributing documents.
3. 🏢 **Use for legitimate purposes**: maintenance, integration, education, research and personal reference.
4. 🤖 **Respect the website.** Review the site's terms of use and `robots.txt` before running the downloader, keep request rates low, and do not run it repeatedly in a loop. Some pages may block automated access.
5. 🔄 **Documents can be outdated.** Manufacturers update manuals. For safety-critical, medical or legal-for-trade use, always confirm against the current official version on the manufacturer's website. 🏥
6. 🩺 **Not medical or metrology advice.** Calibration, certification and legal-for-trade requirements must be handled by qualified professionals. 🧑‍🔧
7. 🚫 **No misuse.** Do not use the material to counterfeit products, forge certifications or mislead customers.
8. 🧯 **Remove on request.** If a rights holder asks for content to be taken down, please honor that request promptly.
9. 🔒 **No personal data.** The archive is meant to contain product documentation only.
10. 📣 **Attribution.** When you share information derived from these documents, credit Detecto as the source.

---

## 🗺️ Roadmap Ideas

These are suggestions for future improvements, not features that exist today. 🚧

- [ ] 📑 Generate a `manifest.csv` listing every file, its source address, size and download date
- [ ] 🔐 Store SHA-256 checksums to detect changed documents
- [ ] 🌐 Support more product groups (chair scales, physician scales, veterinary scales, ...)
- [ ] 🧾 Save original filenames and titles alongside the cleaned names
- [ ] 🗓️ Add a "changed since last run" report
- [ ] 🔤 Auto-extract text for search
- [ ] 🧪 Add automated tests and a GitHub Actions workflow
- [ ] 🏷️ Tag documents by type and language
- [ ] 🤖 Respect `robots.txt` automatically
- [ ] 🐳 Provide a Dockerfile for one-command runs

---

## 🤝 Contributing

Ideas and improvements are welcome! 🎈

- 🐛 Found a bug or a broken link? Open an **issue**.
- ✨ Want to improve the scripts? Send a **pull request**.
- 📝 Spotted an unclear explanation? Documentation fixes are very much appreciated. ❤️

**Suggested workflow:**

```bash
git checkout -b feature/my-improvement
# ...make your changes...
git commit -m "✨ Describe your improvement"
git push origin feature/my-improvement
```

### ✅ Contribution Checklist

- [ ] Code is readable with descriptive names
- [ ] Errors are logged, never silently ignored
- [ ] No third-party dependencies added without discussion (for the Go version)
- [ ] Filenames still follow the sanitization rules
- [ ] README updated if behavior changed

---

## 📄 License

The **source code** in this repository is released under the **MIT License** ⚖️. See the [`LICENSE`](LICENSE) file for details.

The downloaded PDF documents remain the property of their respective owners (see the notes above).

---

<p align="center">
  Made with ❤️, ☕ and a lot of PDFs 📚<br/>
  ⭐ If this repository helped you, consider giving it a star! ⭐
</p>
