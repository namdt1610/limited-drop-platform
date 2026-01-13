// Product Detail Page (Web Component + Alpine)
import { ProductLogic } from "./product.logic.js";

// Định nghĩa custom element: <product-page>
customElements.define(
  "product-page",
  class extends HTMLElement {
    connectedCallback() {
      // Render skeleton và markup chính
      this.innerHTML = /*html*/ `
<div x-data="productPage" class="relative min-h-screen bg-black text-white font-inter selection:bg-white selection:text-black">
  <!-- Background -->
  <mercury-background class="fixed inset-0 z-0 opacity-40"></mercury-background>

  <!-- Header -->
  <header class="reveal py-8 px-6 md:px-12 relative z-10">
    <div class="max-w-screen-2xl mx-auto flex items-center justify-between">
      <a href="#landing" @click.prevent="$dispatch('route', 'landing')" class="group">
        <img src="/imgs/logo-png.webp" alt="DONALD" class="w-20 md:w-24 h-auto object-contain transition-transform duration-500 group-hover:scale-110" />
      </a>
      <div class="flex flex-col items-end">
        <span class="text-[9px] font-mono text-white/30 uppercase tracking-[0.4em] mb-1">Thực thể</span>
        <span class="text-xs font-bold uppercase tracking-widest">Chi tiết Giao thức</span>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="relative z-10 px-6 md:px-12 pb-32">
    <div class="max-w-screen-2xl mx-auto">
      <div x-show="loading" class="reveal flex flex-col items-center py-32 gap-4">
        <div class="w-12 h-12 border-2 border-white/20 border-t-white rounded-full animate-spin"></div>
        <p class="text-[10px] font-mono tracking-[0.3em] text-white/40 uppercase">Đang khởi tạo Thực thể...</p>
      </div>

      <div x-show="error" class="reveal p-8 rounded-3xl bg-red-500/5 border border-red-500/10 text-center">
        <p class="text-xs font-mono text-red-400 uppercase tracking-widest" x-text="error"></p>
      </div>

      <template x-if="product">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-20 items-start">
          <!-- Product Visuals -->
          <div class="lg:col-span-7 reveal" style="transition-delay: 200ms">
            <div class="aspect-square rounded-3xl overflow-hidden bg-zinc-900/50 border border-white/5 backdrop-blur-3xl relative group">
              <img :src="product.thumbnail || (product.images && product.images[0])" alt="Product image" class="w-full h-full object-cover transition-transform duration-1000 group-hover:scale-105" />
              <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-700"></div>
              
              <div class="absolute bottom-8 left-8 right-8 flex items-end justify-between translate-y-4 group-hover:translate-y-0 transition-transform duration-700">
                <div class="space-y-1">
                  <span class="text-[10px] font-mono text-white/40 uppercase tracking-[0.2em]">Mã định danh</span>
                  <p class="text-xs font-mono text-white/80" x-text="'SYMB-' + product.id.toString().padStart(3, '0')"></p>
                </div>
                <div class="text-right">
                  <span class="text-[10px] font-mono text-white/40 uppercase tracking-[0.2em]">Giá trị</span>
                  <p class="text-2xl font-spiky tracking-tighter">
                    <span x-text="new Intl.NumberFormat('vi-VN').format(product.price)"></span>
                    <span class="text-sm text-white/40 ml-1">VND</span>
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Product Info & Actions -->
          <div class="lg:col-span-5 space-y-12 reveal" style="transition-delay: 400ms">
            <div class="space-y-6">
              <div class="space-y-2">
                <h2 class="text-5xl md:text-7xl font-spiky tracking-tighter leading-none uppercase" x-text="product.name"></h2>
                <p class="text-lg text-white/60 font-playfair italic">Một kiệt tác từ bộ sưu tập Cộng sinh.</p>
              </div>

              <div class="py-8 border-y border-white/5">
                <p class="text-sm text-white/40 leading-relaxed" x-text="product.description"></p>
              </div>

              <div class="flex items-center gap-12">
                <div class="space-y-1">
                  <span class="text-[10px] font-mono text-white/30 uppercase tracking-[0.3em]">Trạng thái Kho</span>
                  <div class="flex items-baseline gap-2">
                    <span class="text-2xl font-spiky" x-text="product.stock"></span>
                    <span class="text-xs text-white/40 uppercase tracking-widest">Cá thể khả dụng</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Action -->
            <div class="space-y-6">
              <div class="flex items-center gap-4">
                <div class="flex-1 relative group">
                  <input type="number" x-model="qty" min="1" 
                         class="w-full bg-white/[0.03] border border-white/10 rounded-2xl px-6 py-5 text-lg font-mono text-white focus:outline-none focus:border-white/40 transition-all" />
                  <span class="absolute right-6 top-1/2 -translate-y-1/2 text-[9px] font-mono text-white/20 uppercase tracking-[0.2em]">Số lượng</span>
                </div>
                <button @click="addToCart()" 
                        class="flex-[2] py-6 rounded-full bg-white text-black text-[10px] font-bold tracking-[0.4em] uppercase transition-all duration-500 hover:scale-[1.02]">
                  Thêm vào Giao thức
                </button>
              </div>
              <p class="text-[9px] text-center text-white/30 uppercase tracking-[0.2em]">
                Giao thức Hệ thống: Giao dịch Bảo mật v1.0.4
              </p>
            </div>
          </div>
        </div>
      </template>
    </div>
  </main>

  <!-- Footer -->
  <footer class="reveal py-20 px-6 md:px-12 border-t border-white/5">
    <div class="max-w-screen-2xl mx-auto flex flex-col md:flex-row items-center justify-between gap-8">
      <div class="flex items-center gap-4">
        <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Phiên bản Thực thể</span>
        <span class="text-[9px] font-mono text-white/40 uppercase tracking-[0.4em]">v1.0.4-STABLE</span>
      </div>
      <a href="#landing" @click.prevent="$dispatch('route', 'landing')" class="text-[10px] font-bold uppercase tracking-[0.4em] text-white/40 hover:text-white transition-colors">
        Quay lại Lõi
      </a>
      <div class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">
        © 2026 DONALD CLUB
      </div>
    </div>
  </footer>
</div>
    `;

      // Đăng ký logic tách riêng (guarded)
      if (!Alpine.store("productPageInitialized")) {
        Alpine.data("productPage", ProductLogic);
        Alpine.store("productPageInitialized", true);
      }
    }
  }
);
