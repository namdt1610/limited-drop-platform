import { getProducts } from "../lib/api.js";

export function ProductsLogic() {
  return {
    products: [],
    loading: false,
    error: null,
    page: 1,
    limit: 12,
    total: 0,

    init() {
      this.loadProducts();
      this.initRevealObserver();
    },

    initRevealObserver() {
      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              entry.target.classList.add("reveal-visible");
            }
          });
        },
        { threshold: 0.1 }
      );
      setTimeout(() => {
        document
          .querySelectorAll(".reveal")
          .forEach((el) => observer.observe(el));
      }, 500);
    },

    async loadProducts() {
      this.loading = true;
      this.error = null;
      try {
        const res = await getProducts({ page: this.page, limit: this.limit });
        this.products = res.data || res.products || [];
        this.total = res.total || (res.meta && res.meta.total) || 0;
      } catch (e) {
        this.error = e.message || "Không thể tải danh sách sản phẩm";
      } finally {
        this.loading = false;
      }
    },

    openProduct(p) {
      AlpineRouter.navigate("product", { id: p.id || p.slug });
    },
  };
}
