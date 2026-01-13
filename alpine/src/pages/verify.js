// Verify Page (Web Component + Alpine)
import { VerifyLogic } from "./verify.logic.js";

// Định nghĩa custom element: <verify-page>
customElements.define(
  "verify-page",
  class extends HTMLElement {
    connectedCallback() {
      this.innerHTML = /*html*/ `
<div x-data="verifyPage" class="relative min-h-screen bg-black text-white font-inter selection:bg-white selection:text-black">
  <!-- Background -->
  <mercury-background class="fixed inset-0 z-0 opacity-40"></mercury-background>

  <!-- Header -->
  <header class="reveal py-8 px-6 md:px-12 relative z-10">
    <div class="max-w-screen-2xl mx-auto flex items-center justify-between">
      <a href="#landing" @click.prevent="$dispatch('route', 'landing')" class="group">
        <img src="/imgs/logo-png.webp" alt="DONALD" class="w-20 md:w-24 h-auto object-contain transition-transform duration-500 group-hover:scale-110" />
      </a>
      <div class="flex flex-col items-end">
        <span class="text-[9px] font-mono text-white/30 uppercase tracking-[0.4em] mb-1">Xác thực</span>
        <span class="text-xs font-bold uppercase tracking-widest">Hệ thống SYMBICODE</span>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="relative z-10 px-6 md:px-12 pb-32">
    <div class="max-w-4xl mx-auto">
      <!-- Hero Section -->
      <div class="reveal py-12 md:py-20 space-y-6 text-center">
        <div class="inline-flex items-center gap-3 px-4 py-2 rounded-full bg-purple-500/5 border border-purple-500/10 backdrop-blur-xl">
          <span class="w-1.5 h-1.5 rounded-full bg-purple-500 animate-pulse"></span>
          <span class="text-[10px] font-mono text-purple-400/60 uppercase tracking-[0.3em]">Hệ thống Chống hàng giả</span>
        </div>
        <h1 class="text-5xl md:text-7xl font-spiky tracking-tighter leading-none uppercase">
          XÁC THỰC<br/><span class="text-white/20">SYMBICODE</span>
        </h1>
      </div>

      <div class="reveal p-8 md:p-16 rounded-[2.5rem] bg-zinc-900/50 border border-white/10 backdrop-blur-3xl space-y-12">
        <template x-if="!result">
          <div class="space-y-10">
            <div class="space-y-6">
              <div class="flex items-center gap-4">
                <span class="text-[10px] font-mono text-white/20 uppercase tracking-[0.4em]">01</span>
                <h2 class="text-xs font-bold uppercase tracking-widest">Nhập mã định danh</h2>
              </div>
              
              <div class="relative group">
                <input
                  type="text"
                  x-model="code"
                  @input="code = code.toUpperCase()"
                  placeholder="SYM-XXXXXXXXXXXX"
                  class="w-full bg-white/[0.03] border border-white/10 rounded-2xl px-8 py-6 text-2xl font-mono text-white placeholder:text-white/5 focus:outline-none focus:border-purple-500/40 transition-all"
                  :disabled="verifying"
                />
                <div class="absolute inset-0 rounded-2xl bg-purple-500/5 opacity-0 group-focus-within:opacity-100 pointer-events-none transition-opacity"></div>
              </div>
              
              <p class="text-sm text-white/40 italic leading-relaxed">
                "Mỗi sản phẩm chính hãng sở hữu một mã SYMBICODE duy nhất. Mã này chỉ có thể được kích hoạt một lần để xác nhận quyền sở hữu."
              </p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <button @click="handleVerify" :disabled="verifying || !code" 
                      class="py-6 rounded-full bg-white text-black text-[10px] font-bold tracking-[0.4em] uppercase transition-all duration-500 hover:scale-[1.02] disabled:opacity-20">
                <span x-show="!verifying">Xác thực ngay</span>
                <span x-show="verifying" class="animate-pulse">Đang kiểm tra...</span>
              </button>
              <button @click="handleScan" :disabled="verifying" 
                      class="py-6 rounded-full bg-white/5 border border-white/10 text-white text-[10px] font-bold tracking-[0.4em] uppercase transition-all duration-500 hover:bg-white/10">
                Quét QR Code
              </button>
            </div>
          </div>
        </template>

        <template x-if="result">
          <div class="space-y-12">
            <!-- Success State -->
            <template x-if="result.is_authentic && result.is_first_use">
              <div class="space-y-10 text-center">
                <div class="w-24 h-24 mx-auto bg-green-500/10 rounded-full flex items-center justify-center border border-green-500/20">
                  <svg class="w-10 h-10 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                
                <div class="space-y-4">
                  <h3 class="text-3xl md:text-4xl font-spiky tracking-tighter uppercase text-green-400">XÁC NHẬN VẬT CHỦ</h3>
                  <p class="text-lg text-white/80 leading-relaxed">
                    SYMBIONT ĐÃ ĐƯỢC KÍCH HOẠT VÀO LÚC <span class="font-mono text-green-400" x-text="result.activated_at ? new Date(result.activated_at).toLocaleString('vi-VN') : new Date().toLocaleString('vi-VN')"></span>.
                    <br/>CHÀO MỪNG ĐẾN VỚI DONALD CLUB.
                  </p>
                </div>

                <div class="grid grid-cols-2 gap-4 pt-8 border-t border-white/5">
                  <div class="text-left space-y-1">
                    <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Trạng thái</span>
                    <p class="text-[10px] font-mono text-green-400 uppercase tracking-widest">Đã xác nhận</p>
                  </div>
                  <div class="text-right space-y-1">
                    <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Giao thức</span>
                    <p class="text-[10px] font-mono text-white/60 uppercase tracking-widest">SYMB-v1.0</p>
                  </div>
                </div>
              </div>
            </template>

            <!-- Warning State (Already used or Fake) -->
            <template x-if="!result.is_authentic || !result.is_first_use">
              <div class="space-y-10 text-center">
                <div class="w-24 h-24 mx-auto bg-red-500/10 rounded-full flex items-center justify-center border border-red-500/20 animate-pulse">
                  <svg class="w-10 h-10 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </div>
                
                <div class="space-y-6">
                  <h3 class="text-3xl md:text-4xl font-spiky tracking-tighter uppercase text-red-500">CẢNH BÁO HỆ THỐNG</h3>
                  <div class="p-8 rounded-2xl bg-red-500/5 border border-red-500/10">
                    <p class="text-xl text-red-400 font-bold uppercase tracking-widest leading-relaxed">
                      MÃ NÀY ĐÃ BỊ SỬ DỤNG HOẶC KHÔNG TỒN TẠI.
                      <br/>BẠN ĐANG CẦM TRÊN TAY ĐỐNG RÁC FAKE. VỨT NÓ ĐI.
                    </p>
                  </div>
                </div>

                <template x-if="result.activated_at">
                  <p class="text-[10px] font-mono text-white/20 uppercase tracking-[0.2em]">
                    Lần kích hoạt đầu tiên: <span class="text-white/40" x-text="new Date(result.activated_at).toLocaleString('vi-VN')"></span>
                  </p>
                </template>
              </div>
            </template>

            <div class="pt-8 border-t border-white/5 text-center">
              <button @click="reset" class="text-[10px] font-bold uppercase tracking-[0.4em] text-white/40 hover:text-white transition-colors">
                Kiểm tra mã khác
              </button>
            </div>
          </div>
        </template>
      </div>
    </div>
  </main>

  <!-- Footer -->
  <footer class="reveal py-20 px-6 md:px-12 border-t border-white/5">
    <div class="max-w-screen-2xl mx-auto text-center space-y-4">
      <p class="text-[10px] font-mono text-white/20 uppercase tracking-[0.4em]">
        Donald - Put the world in my cage
      </p>
      <div class="text-[9px] font-mono text-white/10 uppercase tracking-[0.4em]">
        © 2026 DONALD CLUB
      </div>
    </div>
  </footer>
</div>
`;

      // Đăng ký logic tách riêng (guarded)
      if (!Alpine.store("verifyPageInitialized")) {
        Alpine.data("verifyPage", VerifyLogic);
        Alpine.store("verifyPageInitialized", true);
      }
    }
  }
);
