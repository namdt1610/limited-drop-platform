// Policy Page (Web Component + Alpine)
import { PolicyLogic } from "./policy.logic.js";

// Định nghĩa custom element: <policy-page>
customElements.define(
  "policy-page",
  class extends HTMLElement {
    connectedCallback() {
      this.innerHTML = /*html*/ `
<div x-data="policyPage" class="relative min-h-screen bg-black text-white font-inter selection:bg-white selection:text-black">
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
        <span class="text-xs font-bold uppercase tracking-widest">Bảo mật Dữ liệu</span>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="relative z-10 px-6 md:px-12 pb-32">
    <div class="max-w-4xl mx-auto">
      <!-- Hero Section -->
      <div class="reveal py-12 md:py-20 space-y-6 text-center">
        <div class="inline-flex items-center gap-3 px-4 py-2 rounded-full bg-blue-500/5 border border-blue-500/10 backdrop-blur-xl">
          <span class="w-1.5 h-1.5 rounded-full bg-blue-500 animate-pulse"></span>
          <span class="text-[10px] font-mono text-blue-400/60 uppercase tracking-[0.3em]">Giao thức Bảo mật Dữ liệu</span>
        </div>
        <h1 class="text-5xl md:text-7xl font-spiky tracking-tighter leading-none uppercase">
          QUYỀN<br/><span class="text-white/20">RIÊNG TƯ</span>
        </h1>
      </div>

      <div class="reveal p-8 md:p-16 rounded-[2.5rem] bg-zinc-900/50 border border-white/10 backdrop-blur-3xl space-y-16">
        <section class="space-y-6">
          <div class="flex items-center gap-4">
            <span class="text-[10px] font-mono text-white/20 uppercase tracking-[0.4em]">01</span>
            <h2 class="text-xs font-bold uppercase tracking-widest">Lời nói đầu</h2>
          </div>
          <p class="text-lg text-white/60 leading-relaxed font-light pl-8 border-l border-white/5">
            Chúng tôi thu thập dữ liệu tối thiểu cần thiết để xử lý đơn hàng và bảo hành. Sự riêng tư của vật chủ là ưu tiên hàng đầu trong hệ sinh thái Symbiote.
          </p>
        </section>

        <section class="space-y-6">
          <div class="flex items-center gap-4">
            <span class="text-[10px] font-mono text-white/20 uppercase tracking-[0.4em]">02</span>
            <h2 class="text-xs font-bold uppercase tracking-widest">Thông tin thu thập</h2>
          </div>
          <div class="pl-8 border-l border-white/5 space-y-6">
            <p class="text-sm text-white/40 leading-relaxed">
              Chúng tôi chỉ thu thập những thông tin thiết yếu để định danh vật chủ và điều phối vận chuyển:
            </p>
            <ul class="grid grid-cols-1 md:grid-cols-2 gap-4 text-[10px] font-mono uppercase tracking-widest text-white/60">
              <li class="flex items-center gap-3"><span class="w-1 h-1 bg-blue-500 rounded-full"></span> Tên vật chủ</li>
              <li class="flex items-center gap-3"><span class="w-1 h-1 bg-blue-500 rounded-full"></span> Địa chỉ tọa độ</li>
              <li class="flex items-center gap-3"><span class="w-1 h-1 bg-blue-500 rounded-full"></span> Liên lạc (Phone)</li>
              <li class="flex items-center gap-3"><span class="w-1 h-1 bg-blue-500 rounded-full"></span> Email định danh</li>
            </ul>
            <div class="p-6 rounded-2xl bg-blue-500/5 border border-blue-500/10 text-[10px] font-mono text-blue-400/60 uppercase tracking-widest leading-relaxed">
              "Hệ thống không lưu trữ thông tin tài chính. Mọi giao dịch được xử lý qua cổng bảo mật độc lập."
            </div>
          </div>
        </section>

        <section class="space-y-6">
          <div class="flex items-center gap-4">
            <span class="text-[10px] font-mono text-white/20 uppercase tracking-[0.4em]">03</span>
            <h2 class="text-xs font-bold uppercase tracking-widest">Mục đích & Chia sẻ</h2>
          </div>
          <div class="pl-8 border-l border-white/5 grid grid-cols-1 md:grid-cols-2 gap-12">
            <div class="space-y-3">
              <h3 class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Mục đích</h3>
              <p class="text-sm text-white/60 leading-relaxed">Vận chuyển, kích hoạt bảo hành qua SYMBICODE và thông báo trạng thái đơn hàng.</p>
            </div>
            <div class="space-y-3">
              <h3 class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Chia sẻ</h3>
              <p class="text-sm text-white/60 leading-relaxed">Chỉ chia sẻ với đơn vị vận chuyển. Tuyệt đối không bán hoặc cho thuê dữ liệu cho bên thứ ba.</p>
            </div>
          </div>
        </section>

        <div class="pt-12 border-t border-white/5 text-center">
          <button @click="goBack" class="px-12 py-6 rounded-full bg-white text-black text-[10px] font-bold tracking-[0.4em] uppercase transition-all duration-500 hover:scale-[1.02]">
            Xác nhận & Quay lại
          </button>
        </div>
      </div>
    </div>
  </main>

  <!-- Footer -->
  <footer class="reveal py-20 px-6 md:px-12 border-t border-white/5">
    <div class="max-w-screen-2xl mx-auto text-center">
      <div class="text-[9px] font-mono text-white/10 uppercase tracking-[0.4em]">
        © 2026 DONALD CLUB
      </div>
    </div>
  </footer>
</div>
`;

      // Đăng ký logic tách riêng (guarded)
      if (!Alpine.store("policyPageInitialized")) {
        Alpine.data("policyPage", PolicyLogic);
        Alpine.store("policyPageInitialized", true);
      }
    }
  }
);
