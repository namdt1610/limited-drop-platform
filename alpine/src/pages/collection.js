// Collection Page (Web Component + Alpine)
import { CollectionLogic } from "./collection.logic.js";

// Định nghĩa custom element: <collection-page>
customElements.define(
  "collection-page",
  class extends HTMLElement {
    connectedCallback() {
      // Render UI
      this.innerHTML = /*html*/ `
<div x-data="collectionPage" class="relative min-h-screen bg-black text-white font-inter selection:bg-white selection:text-black">
  <!-- Background -->
  <mercury-background class="fixed inset-0 z-0 opacity-40"></mercury-background>

  <!-- Header -->
  <header class="reveal py-8 px-6 md:px-12 relative z-10">
    <div class="max-w-screen-2xl mx-auto flex items-center justify-between">
      <a href="#landing" @click.prevent="$dispatch('route', 'landing')" class="group">
        <img src="/imgs/logo-png.webp" alt="DONALD" class="w-20 md:w-24 h-auto object-contain transition-transform duration-500 group-hover:scale-110" />
      </a>
      <div class="flex flex-col items-end">
        <span class="text-[9px] font-mono text-white/30 uppercase tracking-[0.4em] mb-1">Giao thức</span>
        <span class="text-xs font-bold uppercase tracking-widest">Toàn bộ Bộ sưu tập</span>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="relative z-10 px-6 md:px-12 pb-32">
    <div class="max-w-screen-2xl mx-auto">
      <!-- Hero Section -->
      <div class="reveal py-20 md:py-32 space-y-8 text-center">
        <div class="inline-flex items-center gap-3 px-4 py-2 rounded-full bg-white/5 border border-white/10 backdrop-blur-xl">
          <span class="w-1.5 h-1.5 rounded-full bg-green-500 animate-pulse"></span>
          <span class="text-[10px] font-mono text-white/60 uppercase tracking-[0.3em]">Hệ thống Trực tuyến</span>
        </div>
        <h1 class="text-6xl md:text-9xl font-spiky tracking-tighter leading-none uppercase">
          CÁC<br/><span class="text-white/20">CÁ THỂ CỘNG SINH</span>
        </h1>
        <p class="max-w-2xl mx-auto text-lg md:text-xl text-white/40 font-playfair italic">
          Sự giao thoa giữa vẻ đẹp sinh học và độ chính xác cơ khí. Mỗi cá thể là một giao thức duy nhất trong hệ sinh thái Symbiosis.
        </p>
      </div>

      <!-- Products Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 md:gap-12">
        <!-- Active Product -->
        <div class="reveal group" style="transition-delay: 200ms">
          <a href="#drop" @click.prevent="$dispatch('route', { page: 'drop', id: '1' })" class="block space-y-6">
            <div class="aspect-square rounded-3xl overflow-hidden bg-zinc-900/50 border border-white/5 backdrop-blur-3xl relative">
              <img src="/imgs/collection-main.webp" alt="Symbiote Chrome #001" class="w-full h-full object-cover transition-transform duration-1000 group-hover:scale-105" />
              <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-700"></div>
              <div class="absolute top-6 right-6">
                <span class="px-3 py-1 rounded-full bg-green-500/10 border border-green-500/20 text-[9px] font-bold text-green-500 uppercase tracking-widest">Khả dụng</span>
              </div>
              <div class="absolute bottom-8 left-8 translate-y-4 group-hover:translate-y-0 transition-transform duration-700">
                <span class="text-[10px] font-mono text-white/40 uppercase tracking-[0.2em]">Giao thức #001</span>
                <p class="text-xl font-spiky tracking-tighter uppercase">Cá thể Alpha</p>
              </div>
            </div>
          </a>
        </div>

        <!-- Locked Products -->
        <template x-for="num in [2, 3, 4, 5, 6]" :key="num">
          <div class="reveal group opacity-40 grayscale hover:grayscale-0 transition-all duration-700" :style="'transition-delay: ' + (num * 100) + 'ms'">
            <div class="aspect-square rounded-3xl overflow-hidden bg-white/[0.02] border border-white/5 backdrop-blur-sm relative flex items-center justify-center">
              <div class="text-center space-y-4">
                <span class="text-6xl font-spiky text-white/5 select-none" x-text="'#' + num.toString().padStart(3, '0')"></span>
                <div class="flex flex-col items-center gap-2">
                  <div class="w-8 h-[1px] bg-white/10"></div>
                  <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Giao thức bị Khóa</span>
                </div>
              </div>
              
              <!-- Hover Info -->
              <div class="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-500 bg-black/40 backdrop-blur-md">
                <p class="text-[10px] font-mono text-white/60 uppercase tracking-[0.3em]">Đang chờ Kích hoạt</p>
              </div>
            </div>
            <div class="mt-6 space-y-1">
              <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.2em]" x-text="'SYMB-00' + num"></span>
              <p class="text-sm font-bold uppercase tracking-widest text-white/40">Thực thể Chưa xác định</p>
            </div>
          </div>
        </template>
      </div>
    </div>
  </main>

  <!-- Footer -->
  <footer class="reveal py-20 px-6 md:px-12 border-t border-white/5">
    <div class="max-w-screen-2xl mx-auto flex flex-col md:flex-row items-center justify-between gap-8">
      <div class="flex items-center gap-4">
        <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Phiên bản Hệ thống</span>
        <span class="text-[9px] font-mono text-white/40 uppercase tracking-[0.4em]">v1.0.4-STABLE</span>
      </div>
      <a href="#" @click.prevent="goHome" class="text-[10px] font-bold uppercase tracking-[0.4em] text-white/40 hover:text-white transition-colors">
        Quay lại Lõi
      </a>
      <div class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">
        © 2026 DONALD CLUB
      </div>
    </div>
  </footer>
</div>`;

      // Đăng ký logic tách riêng (guarded)
      if (!Alpine.store("collectionPageInitialized")) {
        Alpine.data("collectionPage", CollectionLogic);
        Alpine.store("collectionPageInitialized", true);
      }
    }
  }
);
