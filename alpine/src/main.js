import Alpine from "alpinejs";
import "./../index.css";

// 0. Register background component first
import "./components/landing/mercury-background.js";

// 1. Register Teaser / Landing components early so they exist before pages render
import "./components/landing/teaser-footer.js";
import "./components/landing/teaser-story-section.js";
import "./components/landing/teaser-hero-section.js";
import "./components/landing/system-stories-section.js";
import "./components/landing/watch-carousel.js";
import "./components/landing/cinematic-video-section.js";
import "./components/landing/waitlist-form.js";

// 2. Đăng ký các Web Components (pages) — after primitives & teasers
import "./pages/landing.js";
import "./pages/drop.js";
import "./pages/product.js";
import "./pages/products.js";
import "./pages/collection.js";
import "./pages/payment-success.js";
import "./pages/payment-cancel.js";
import "./pages/terms.js";
import "./pages/policy.js";
import "./pages/track-order.js";
import "./pages/verify.js";

// Ensure custom elements are block-level to avoid layout issues
const style = document.createElement("style");
style.textContent = `
  landing-page, drop-page, product-page, products-page, collection-page, 
  checkout-page, payment-success-page, payment-cancel-page, terms-page, 
  policy-page, track-order-page, verify-page {
    display: block;
  }
`;
document.head.appendChild(style);

window.getHashParams = () => {
  const hash = window.location.hash.replace("#", "");
  const [page, queryString] = hash.split("?");
  const params = {};
  if (queryString) {
    const urlParams = new URLSearchParams(queryString);
    urlParams.forEach((value, key) => {
      params[key] = value;
    });
  }
  return params;
};

document.addEventListener("alpine:init", () => {
  // 2. Định nghĩa Router Handler (DDO)
  Alpine.data("router", () => ({
    path: "landing", // Default page

    init() {
      const handleHash = () => {
        const hash = window.location.hash.replace("#", "") || "landing";
        const [page, queryString] = hash.split("?");
        const params = {};
        if (queryString) {
          const urlParams = new URLSearchParams(queryString);
          urlParams.forEach((value, key) => {
            params[key] = value;
          });
        }
        this.path = page;
        window.dispatchEvent(
          new CustomEvent("route-changed", {
            detail: { route: page, params: params },
          })
        );
      };

      window.addEventListener("hashchange", handleHash);
      handleHash();
    },

    navigate(detail) {
      const page = typeof detail === "string" ? detail : detail.page;
      let hash = page;
      let params = {};

      if (typeof detail === "object") {
        params = { ...detail };
        delete params.page;
        const query = new URLSearchParams(params).toString();
        if (query) hash += "?" + query;
      }

      // Setting hash will trigger hashchange event which calls handleHash
      window.location.hash = hash;
    },
  }));
});

window.Alpine = Alpine;
Alpine.start();
