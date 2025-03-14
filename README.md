# ✂️ **GoShort** - Fast and customizable URL shortener

![Docker Image Version](https://img.shields.io/docker/v/petrakisg/goshort?sort=semver&label=Docker%20Image%20Version&logo=docker)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kek-Sec/GoShort)


### Demo: [https://x.yup.gr](http://x.yup.gr)
### DockerHub Image: [petrakisg/goshort](https://hub.docker.com/r/petrakisg/goshort)

GoShort is a fast and customizable URL shortener built with Go, Svelte and TailwindCSS. It is designed to be self-hosted and easy to deploy.

![GoShort](web/static/banner.png)

---

## 📋 **Table of Contents**

1. [Features](#-features)  
2. [Installation](#-installation)  
3. [Customization](#-customization)
4. [Contributing](#-contributing)
5. [License](#-license)
6. [Security](#-security)
---

## 🚀 **Features**

- **Fast**: GoShort is built with Go and is blazing fast.
- **Customizable**: GoShort is built with Svelte and TailwindCSS, making it easy to customize.
- **Self-hosted**: You own your data and can deploy GoShort on your own server.
- **Custom URLs**: You can set custom URLs for your short links.
- **Expiration**: You can set expiration for your short links.
- **White-labeling**: Customize branding elements without rebuilding the Docker image.

---

## 🛠️ **Installation**

> Checkout docker-compose.prod.yml for a sample production setup.

---

## 🎨 **Customization**

GoShort supports customizing the branding and appearance through environment variables, making it easy to white-label without rebuilding the Docker image.

### Available Customization Options

You can customize the following aspects of the UI by setting these environment variables:

| Environment Variable | Description | Default Value |
|---------------------|-------------|---------------|
| BRAND_TITLE | Browser tab title | GoShort - URL Shortener |
| BRAND_DESCRIPTION | Meta description for SEO | GoShort is a powerful and user-friendly URL shortener... |
| BRAND_KEYWORDS | Meta keywords for SEO | URL shortener, GoShort, link management... |
| BRAND_AUTHOR | Author meta tag | GoShort Team |
| BRAND_THEME_COLOR | Browser theme color | #4caf50 |
| BRAND_LOGO_TEXT | Text logo displayed in the header | GoShort |
| BRAND_PRIMARY_COLOR | Main accent color (buttons, links) | #3b82f6 |
| BRAND_SECONDARY_COLOR | Secondary accent color | #10b981 |
| BRAND_HEADER_TITLE | Main heading on the page | GoShort - URL Shortener |
| BRAND_FOOTER_TEXT | Text shown in the footer | View the project on |
| BRAND_FOOTER_LINK | URL for the footer link | https://github.com/kek-Sec/GoShort |

### Usage Example

Here's how to customize the branding in your docker-compose file:

```yaml
services:
  goshort:
    image: petrakisg/goshort:1.0.1
    environment:
      # Database configuration
      DATABASE_URL: postgres://user:password@db:5432/goshort
      
      # Branding customization
      BRAND_TITLE: "MyCompany URL Shortener"
      BRAND_LOGO_TEXT: "MyShort"
      BRAND_PRIMARY_COLOR: "#ff5722"
      BRAND_SECONDARY_COLOR: "#2196f3"
      BRAND_HEADER_TITLE: "MyCompany Link Shortener"
      BRAND_FOOTER_TEXT: "Powered by"
      BRAND_FOOTER_LINK: "https://mycompany.com"
```

---

## 🤝 **Contributing**

1. Fork the repository.  
2. Create a new branch: `git checkout -b my-feature-branch`  
3. Make your changes and add tests.  
4. Submit a pull request.  

---

## 📄 **License**

GoShort is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.

---

## 🔒 **Security**

We take security seriously and appreciate your efforts to responsibly disclose vulnerabilities. Checkout [SECURITY.md](SECURITY.md) for more information.