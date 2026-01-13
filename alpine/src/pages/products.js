// Products Listing Page (Web Component + Alpine)
import { ProductsLogic } from "./products.logic.js";

// Định nghĩa custom element: <products-page>
customElements.define(
  "products-page",
  class extends HTMLElement {
    connectedCallback() {
      this.innerHTML = /*html*/ `
<div x-data="productsPage" class="relative min-h-screen bg-black text-white font-inter selection:bg-white selection:text-black">
  <!-- Background -->
  <mercury-background class="fixed inset-0 z-0 opacity-40"></mercury-background>

  <!-- Header -->
  <header class="reveal py-8 px-6 md:px-12 relative z-10">
    <div class="max-w-screen-2xl mx-auto flex items-center justify-between">
      <a href="#landing" @click.prevent="$dispatch('route', 'landing')" class="group">
        <img src="/imgs/logo-png.webp" alt="DONALD" class="w-20 md:w-24 h-auto object-contain transition-transform duration-500 group-hover:scale-110" />
      </a>
      <div class="flex flex-col items-end">
        <span class="text-[9px] font-mono text-white/30 uppercase tracking-[0.4em] mb-1">Lưu trữ</span>
        <span class="text-xs font-bold uppercase tracking-widest">Cơ sở dữ liệu Sản phẩm</span>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="relative z-10 px-6 md:px-12 pb-32">
    <div class="max-w-screen-2xl mx-auto">
      <!-- Hero Section -->
      <div class="reveal py-20 md:py-32 space-y-8">
        <div class="inline-flex items-center gap-3 px-4 py-2 rounded-full bg-white/5 border border-white/10 backdrop-blur-xl">
          <span class="w-1.5 h-1.5 rounded-full bg-white/20"></span>
          <span class="text-[10px] font-mono text-white/60 uppercase tracking-[0.3em]">Đang truy cập Cơ sở dữ liệu...</span>
        </div>
        <h1 class="text-6xl md:text-9xl font-spiky tracking-tighter leading-none uppercase">
          KHO<br/><span class="text-white/20">LƯU TRỮ</span>
        </h1>
      </div>

      <div x-show="loading" class="reveal flex flex-col items-center py-20 gap-4">
        <div class="w-12 h-12 border-2 border-white/20 border-t-white rounded-full animate-spin"></div>
        <p class="text-[10px] font-mono tracking-[0.3em] text-white/40 uppercase">Đang giải mã tệp tin...</p>
      </div>

      <div x-show="error" class="reveal p-8 rounded-3xl bg-red-500/5 border border-red-500/10 text-center">
        <p class="text-xs font-mono text-red-400 uppercase tracking-widest" x-text="error"></p>
      </div>

      <!-- Products Grid -->
      <div x-show="!loading && !error" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 md:gap-10">
        <template x-for="(p, index) in products" :key="p.id">
          <div class="reveal group" :style="'transition-delay: ' + (index * 100) + 'ms'">
            <div @click="openProduct(p)" class="cursor-pointer space-y-6">
              <div class="aspect-[4/5] rounded-3xl overflow-hidden bg-zinc-900/50 border border-white/5 backdrop-blur-3xl relative">
                <img :src="p.thumbnail || (p.images && p.images[0])" alt="Product image" class="w-full h-full object-cover transition-transform duration-1000 group-hover:scale-105" />
                <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-700"></div>
                
                <div class="absolute bottom-6 left-6 right-6 translate-y-4 group-hover:translate-y-0 transition-transform duration-700 flex items-end justify-between">
                  <div class="space-y-1">
                    <span class="text-[9px] font-mono text-white/40 uppercase tracking-[0.2em]">Mã định danh</span>
                    <p class="text-[10px] font-mono text-white/80" x-text="'SYMB-' + p.id.toString().padStart(3, '0')"></p>
                  </div>
                  <div class="text-right">
                    <span class="text-[9px] font-mono text-white/40 uppercase tracking-[0.2em]">Giá trị</span>
                    <p class="text-lg font-spiky tracking-tighter" x-text="new Intl.NumberFormat('vi-VN').format(p.price) + ' VND'"></p>
                  </div>
                </div>
              </div>
              <div class="space-y-1 px-2">
                <h3 class="text-xs font-bold uppercase tracking-widest group-hover:text-white transition-colors" x-text="p.name"></h3>
                <p class="text-[10px] font-mono text-white/20 uppercase tracking-[0.2em]">Giao thức đã Xác minh</p>
              </div>
            </div>
          </div>
        </template>
      </div>

      <!-- Pagination -->
      <div x-show="!loading && products.length > 0" class="reveal mt-20 flex items-center justify-center gap-8">
        <button @click="page = page - 1; loadProducts()" :disabled="page <= 1" 
                class="text-[10px] font-bold uppercase tracking-[0.4em] text-white/20 hover:text-white disabled:opacity-0 transition-all">
          Trang_Trước
        </button>
        <div class="flex items-center gap-4">
          <span class="w-8 h-[1px] bg-white/10"></span>
          <span class="text-[10px] font-mono text-white/40" x-text="page"></span>
          <span class="w-8 h-[1px] bg-white/10"></span>
        </div>
        <button @click="page = page + 1; loadProducts()" :disabled="products.length < limit" 
                class="text-[10px] font-bold uppercase tracking-[0.4em] text-white/20 hover:text-white disabled:opacity-0 transition-all">
          Trang_Sau
        </button>
      </div>
    </div>
  </main>

  <!-- Footer -->
  <footer class="reveal py-20 px-6 md:px-12 border-t border-white/5">
    <div class="max-w-screen-2xl mx-auto flex flex-col md:flex-row items-center justify-between gap-8">
      <div class="flex items-center gap-4">
        <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Phiên bản Cơ sở dữ liệu</span>
        <span class="text-[9px] font-mono text-white/40 uppercase tracking-[0.4em]">DB-v2.1.0</span>
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
      if (!Alpine.store("productsPageInitialized")) {
        Alpine.data("productsPage", ProductsLogic);
        Alpine.store("productsPageInitialized", true);
      }
    }
  }
);
