/**
 * @type {import('gatsby').GatsbyConfig}
 */
require("dotenv").config({
  path: `.env`,
})

module.exports = {
  pathPrefix: `/kubernetes-security-assessment`,
  siteMetadata: {
    title: `Kubernetes Security Assessment`,
    siteUrl: `https://pavel-semenov-1.github.io/kubernetes-security-assessment`
  },
  plugins: [{
    resolve: 'gatsby-plugin-google-analytics',
    options: {
      "trackingId": "G-PE16R6WHTG"
    }
  }, "gatsby-transformer-remark", {
    resolve: 'gatsby-source-filesystem',
    options: {
      "name": "pages",
      "path": "./src/pages/"
    },
    __key: "pages"
  }]
};