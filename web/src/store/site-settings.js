let SiteSettings = {
  siteName: {
    type: "v-text-field",
    disabled: false,
    value: "Gocipe",
    attr: {
      icon: "short_text",
      label: "Site Name"
    }
  },
  logo: {
    type: "v-text-field",
    disabled: false,
    value: "",
    attr: {
      icon: "short_text",
      label: "Logo path"
    }
  }
};

module.exports = SiteSettings;
