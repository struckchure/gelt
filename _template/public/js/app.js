document.addEventListener("alpine:init", () => {
  Alpine.data("root", () => ({
    async loadContent(page) {
      page = window.location.origin + "/" + "_" + page;
      if (!page.endsWith("/")) {
        page += "/";
      }

      root = document.getElementById("root");
      try {
        const response = await fetch(page);
        if (!response.ok) throw new Error("Page not found");
        root.innerHTML = await response.text();
      } catch (error) {
        root.innerHTML = `<p style="color: red;">Page not found.</p>`;
      }
    },
    async init() {
      const url = window.location.href.replace(window.location.origin, "");
      await this.loadContent(url);

      document.body.addEventListener("click", async (event) => {
        const link = event.target.closest("a"); // Ensure the click is on a link
        if (link && link.href && link.origin === window.location.origin) {
          event.preventDefault(); // Prevent full-page reload

          const path = link.getAttribute("href");
          history.pushState({ path }, "", path);

          await this.loadContent(path); // Replace this with your SPA's content loader function
        }
      });

      window.addEventListener("popstate", async (event) => {
        const page = event.state?.page || "";
        await this.loadContent(page);
      });
    },
  }));
});
