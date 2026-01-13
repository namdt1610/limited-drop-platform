import { DropLogic } from "./drop.logic.js";
import { HeaderComponent } from "../components/drop/header.js";
import { ProductSection } from "../components/drop/product-section.js";
import { CheckoutModal } from "../components/drop/checkout-modal.js";

customElements.define(
  "drop-page",
  class extends HTMLElement {
    connectedCallback() {
      // Register Alpine data when element is connected (Alpine is now available)
      if (!window.dropPageDataRegistered) {
        Alpine.data("dropPage", DropLogic);
        window.dropPageDataRegistered = true;
      }

      this.innerHTML = /*html*/ `
<div x-data="dropPage" x-init="init()" class="relative min-h-screen bg-black text-white font-inter selection:bg-white selection:text-black">
  <!-- Background -->
  <mercury-background class="fixed inset-0 z-0 opacity-40"></mercury-background>

  <!-- Loading Skeleton -->
  <div x-show="isLoading" class="fixed inset-0 z-[100] flex items-center justify-center bg-black">
    <div class="flex flex-col items-center gap-4">
      <div class="w-12 h-12 border-2 border-white/20 border-t-white rounded-full animate-spin"></div>
      <p class="text-[10px] font-mono tracking-[0.3em] text-white/40 uppercase">Khởi tạo Giao thức...</p>
    </div>
  </div>

  <!-- Drop Page Content -->
  <div x-show="!isLoading" class="relative z-10">
    ${HeaderComponent()}

    <!-- Main Content -->
    <main class="px-6 md:px-12 pb-32">
      <div class="max-w-screen-2xl mx-auto">
        ${ProductSection()}
      </div>
    </main>
  </div>

  ${CheckoutModal()}
</div>`;
    }
  }
);
