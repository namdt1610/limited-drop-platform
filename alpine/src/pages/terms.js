// Terms Page (Web Component + Alpine)
import { TermsLogic } from "./terms.logic.js";

// Định nghĩa custom element: <terms-page>
customElements.define(
  "terms-page",
  class extends HTMLElement {
    connectedCallback() {
      this.innerHTML = /*html*/ `
<div x-data="termsPage" class="relative min-h-screen bg-black text-white font-inter selection:bg-white selection:text-black">
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
        <span class="text-xs font-bold uppercase tracking-widest">Điều khoản Hệ thống</span>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="relative z-10 px-6 md:px-12 pb-32">
    <div class="max-w-4xl mx-auto">
      <!-- Hero Section -->
      <div class="reveal py-12 md:py-20 space-y-6 text-center">
        <div class="inline-flex items-center gap-3 px-4 py-2 rounded-full bg-red-500/5 border border-red-500/10 backdrop-blur-xl">
          <span class="w-1.5 h-1.5 rounded-full bg-red-500 animate-pulse"></span>
          <span class="text-[10px] font-mono text-red-400/60 uppercase tracking-[0.3em]">Thỏa thuận Cộng sinh</span>
        </div>
        <h1 class="text-5xl md:text-7xl font-spiky tracking-tighter leading-none uppercase">
          ĐIỀU KHOẢN<br/><span class="text-white/20">HỆ THỐNG</span>
        </h1>
      </div>

      <div class="reveal p-8 md:p-16 rounded-[2.5rem] bg-zinc-900/50 border border-white/10 backdrop-blur-3xl space-y-16">
        <div class="p-6 rounded-2xl bg-red-500/5 border border-red-500/10 text-center">
          <p class="text-[10px] font-mono text-red-400/80 uppercase tracking-[0.2em]">
            "Khi thanh toán thành công, bạn chấp nhận các điều khoản này."
          </p>
        </div>

        <section class="space-y-6">
          <div class="flex items-center gap-4">
            <span class="text-[10px] font-mono text-white/20 uppercase tracking-[0.4em]">01</span>
            <h2 class="text-xs font-bold uppercase tracking-widest">Định nghĩa</h2>
          </div>
          <div class="pl-8 border-l border-white/5 grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="p-4 rounded-xl bg-white/5 border border-white/5 space-y-1">
              <span class="text-[9px] font-mono text-red-500 uppercase tracking-widest">/ CHÚNG TÔI</span>
              <p class="text-xs font-bold text-white/80">Donald Vibe</p>
            </div>
            <div class="p-4 rounded-xl bg-white/5 border border-white/5 space-y-1">
              <span class="text-[9px] font-mono text-red-500 uppercase tracking-widest">/ BẠN</span>
              <p class="text-xs font-bold text-white/80">Vật chủ (Buyer)</p>
            </div>
            <div class="p-4 rounded-xl bg-white/5 border border-white/5 space-y-1">
              <span class="text-[9px] font-mono text-red-500 uppercase tracking-widest">/ SẢN PHẨM</span>
              <p class="text-xs font-bold text-white/80">Symbiote Gear</p>
            </div>
          </div>
        </section>

        <section class="space-y-6">
          <div class="flex items-center gap-4">
            <span class="text-[10px] font-mono text-white/20 uppercase tracking-[0.4em]">02</span>
            <h2 class="text-xs font-bold uppercase tracking-widest">Quy tắc Giao dịch</h2>
          </div>
          <ul class="pl-8 border-l border-white/5 space-y-6">
            <li class="flex gap-4">
              <span class="text-red-500 font-mono text-xs">!</span>
              <p class="text-sm text-white/60 leading-relaxed">
                <strong class="text-white uppercase tracking-widest text-[10px]">Không giữ chỗ:</strong> Thêm vào giỏ hàng không đảm bảo sở hữu. Chỉ khi thanh toán thành công.
              </p>
            </li>
            <li class="flex gap-4">
              <span class="text-red-500 font-mono text-xs">!</span>
              <p class="text-sm text-white/60 leading-relaxed">
                <strong class="text-white uppercase tracking-widest text-[10px]">Tốc độ quyết định:</strong> Hết hàng là hết. Ai nhanh tay người đó có.
              </p>
            </li>
            <li class="flex gap-4">
              <span class="text-red-500 font-mono text-xs">!</span>
              <p class="text-sm text-white/60 leading-relaxed">
                <strong class="text-white uppercase tracking-widest text-[10px]">Thông tin chính xác:</strong> Sai địa chỉ/số điện thoại → đơn hàng bị hủy.
              </p>
            </li>
          </ul>
        </section>

        <section class="space-y-6">
          <div class="flex items-center gap-4">
            <span class="text-[10px] font-mono text-white/20 uppercase tracking-[0.4em]">03</span>
            <h2 class="text-xs font-bold uppercase tracking-widest">Vận chuyển & Niêm phong</h2>
          </div>
          <div class="pl-8 border-l border-white/5 space-y-8">
            <p class="text-sm text-white/60 leading-relaxed">Giao hàng trong 3-5 ngày làm việc. Tất cả gói hàng có tem niêm phong.</p>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="p-6 rounded-2xl bg-white/5 border border-white/5 space-y-3">
                <h3 class="text-[9px] font-mono text-red-400 uppercase tracking-[0.4em]">Cảnh báo</h3>
                <p class="text-xs text-white/40 italic leading-relaxed">Tem rách hoặc có dấu hiệu mở → từ chối nhận ngay lập tức.</p>
              </div>
              <div class="p-6 rounded-2xl bg-white/5 border border-white/5 space-y-3">
                <h3 class="text-[9px] font-mono text-red-400 uppercase tracking-[0.4em]">Yêu cầu</h3>
                <p class="text-xs text-white/40 italic leading-relaxed">Bắt buộc quay video mở hộp. Không có video → không giải quyết khiếu nại.</p>
              </div>
            </div>
          </div>
        </section>

        <section class="space-y-6">
          <div class="flex items-center gap-4">
            <span class="text-[10px] font-mono text-white/20 uppercase tracking-[0.4em]">04</span>
            <h2 class="text-xs font-bold uppercase tracking-widest">Đổi trả & Bảo hành</h2>
          </div>
          <div class="pl-8 border-l border-white/5 space-y-6">
            <p class="text-sm text-white/60 leading-relaxed">Không hoàn tiền vì lý do chủ quan. Chỉ đổi mới nếu sản phẩm lỗi kỹ thuật trong 7 ngày đầu.</p>
            <div class="p-6 rounded-2xl bg-blue-500/5 border border-blue-500/10">
              <p class="text-[10px] font-mono text-blue-400/60 uppercase tracking-widest leading-relaxed">
                Mỗi sản phẩm có SYMBICODE riêng. Mất mã → mất quyền bảo hành.
              </p>
            </div>
          </div>
        </section>

        <div class="pt-12 border-t border-white/5 text-center">
          <button @click="goBack" class="px-12 py-6 rounded-full bg-white text-black text-[10px] font-bold tracking-[0.4em] uppercase transition-all duration-500 hover:scale-[1.02]">
            Chấp nhận & Quay lại
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
      if (!Alpine.store("termsPageInitialized")) {
        Alpine.data("termsPage", TermsLogic);
        Alpine.store("termsPageInitialized", true);
      }
    }
  }
);
