import { getProduct } from "../lib/api.js";
import { initRevealObserver } from "../lib/dom-utils.js";

export function ProductLogic() {
  return {
    product: null,
    loading: false,
    error: null,
    qty: 1,

    init() {
      this.loadProductFromParams();
      this.initRevealObserver();

      window.addEventListener("route-changed", (event) => {
        if (event.detail.route === "product") {
          this.loadProductFromParams();
        }
      });
    },

    initRevealObserver() {
      initRevealObserver();
    },

    async loadProductFromParams() {
      const params = window.getHashParams();
      const id = params.id;
      if (!id) return;
      this.loading = true;
      try {
        this.product = await getProduct(id);
      } catch (e) {
        this.error = e.message || "Không thể tải sản phẩm";
      } finally {
        this.loading = false;
      }
    },

    addToCart() {
      if (!this.product) return;
      const cart = JSON.parse(localStorage.getItem("cart") || "[]");
      const item = cart.find((i) => i.id === this.product.id);
      if (item) item.qty = item.qty + this.qty;
      else
        cart.push({
          id: this.product.id,
          name: this.product.name,
          price: this.product.price,
          qty: this.qty,
        });
      localStorage.setItem("cart", JSON.stringify(cart));
      alert("Đã thêm vào giỏ hàng");
    },
  };
}
