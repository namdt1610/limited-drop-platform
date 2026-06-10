import { ShowcaseLogic } from "./showcase.logic.js";

customElements.define(
    "showcase-page",
    class extends HTMLElement {
        connectedCallback() {
            this.innerHTML = /*html*/ `
<div x-data="showcasePage" class="relative min-h-screen bg-black text-white font-inter selection:bg-white selection:text-black">
  <!-- Background -->
  <mercury-background class="fixed inset-0 z-0 opacity-40"></mercury-background>

  <!-- Header -->
  <header class="reveal py-8 px-6 md:px-12 relative z-10">
    <div class="max-w-screen-2xl mx-auto flex items-center justify-between">
      <a href="#landing" @click.prevent="$dispatch('route', 'landing')" class="group">
        <img src="/imgs/logo-png.webp" alt="DONALD" class="w-20 md:w-24 h-auto object-contain transition-transform duration-500 group-hover:scale-110" />
      </a>
      <div class="flex flex-col items-end">
        <span class="text-[9px] font-mono text-white/30 uppercase tracking-[0.4em] mb-1">Trưng bày</span>
        <span class="text-xs font-bold uppercase tracking-widest">Bộ sưu tập Danh dự</span>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="relative z-10 px-6 md:px-12 pb-32">
    <div class="max-w-screen-2xl mx-auto">
      <!-- Hero Section -->
      <div class="reveal py-20 md:py-32 space-y-8">
        <div class="inline-flex items-center gap-3 px-4 py-2 rounded-full bg-white/5 border border-white/10 backdrop-blur-xl">
          <span class="w-1.5 h-1.5 rounded-full bg-cyan-500 animate-pulse"></span>
          <span class="text-[10px] font-mono text-white/60 uppercase tracking-[0.3em]">Hệ thống Showcase Đã Kích hoạt</span>
        </div>
        <h1 class="text-6xl md:text-9xl font-spiky tracking-tighter leading-none uppercase">
          TRIỂN<br/><span class="text-white/20">LÃM MẪU</span>
        </h1>
        <p class="max-w-xl text-sm md:text-base text-white/40 leading-relaxed font-light">
          Khám phá những thiết kế tiêu biểu nhất từ DONALD CLUB. Những sản phẩm này thể hiện tinh thần <span class="text-white">Symbiosis Chrome</span> — nơi công nghệ và thời trang hòa quyện.
        </p>
      </div>

      <!-- Products Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 md:gap-10">
        <template x-for="(p, index) in products" :key="p.id">
          <div class="reveal group" :style="'transition-delay: ' + (index * 100) + 'ms'">
            <div @click="openProduct(p)" class="cursor-pointer space-y-6">
              <div class="aspect-4/5 rounded-3xl overflow-hidden bg-zinc-900/50 border border-white/5 backdrop-blur-3xl relative">
                <img :src="p.thumbnail" alt="Product image" class="w-full h-full object-cover transition-transform duration-1000 group-hover:scale-110" loading="lazy" />
                
                <!-- Category Tag -->
                <div class="absolute top-6 left-6">
                    <span class="px-3 py-1 bg-black/60 backdrop-blur-md rounded-full text-[8px] font-mono uppercase tracking-widest border border-white/10" x-text="p.category"></span>
                </div>

                <div class="absolute inset-0 bg-linear-to-t from-black/90 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-700"></div>
                
                <div class="absolute bottom-8 left-8 right-8 translate-y-4 group-hover:translate-y-0 transition-transform duration-700 flex items-end justify-between">
                  <div class="space-y-1">
                    <span class="text-[9px] font-mono text-white/40 uppercase tracking-[0.2em]">Sản phẩm</span>
                    <p class="text-[10px] font-mono text-white/80" x-text="p.name"></p>
                  </div>
                  <div class="text-right">
                    <span class="text-[9px] font-mono text-white/40 uppercase tracking-[0.2em]">Giá</span>
                    <p class="text-base font-spiky tracking-tighter" x-text="new Intl.NumberFormat('vi-VN').format(p.price) + ' đ'"></p>
                  </div>
                </div>
              </div>
              
              <!-- Subtle label below -->
              <div class="flex items-center justify-between px-2">
                <span class="text-[10px] font-mono text-white/10 uppercase tracking-[0.3em]" x-text="'MOCK-ID: ' + p.id"></span>
                <div class="w-4 h-px bg-white/10"></div>
              </div>
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
        <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Trạng thái Hệ thống</span>
        <span class="text-[9px] font-mono text-cyan-500 uppercase tracking-[0.4em]">MÔ PHỎNG_ONLINE</span>
      </div>
      <a href="#landing" @click.prevent="$dispatch('route', 'landing')" class="text-[10px] font-bold uppercase tracking-[0.4em] text-white/40 hover:text-white transition-colors">
        Trở về Trung tâm
      </a>
      <div class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">
        © 2026 DONALD CLUB SHOWCASE
      </div>
    </div>
  </footer>
</div>
      `;

            if (!Alpine.store("showcasePageInitialized")) {
                Alpine.data("showcasePage", ShowcaseLogic);
                Alpine.store("showcasePageInitialized", true);
            }
        }
    }
);
