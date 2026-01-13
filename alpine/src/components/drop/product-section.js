export function ProductSection() {
  return /*html*/ `
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-20 items-start">
      
      <!-- Product Visuals -->
      <div class="lg:col-span-7 reveal" style="transition-delay: 200ms">
        <div class="relative group">
          <div class="aspect-[4/5] md:aspect-square rounded-3xl overflow-hidden bg-zinc-900/50 border border-white/5 backdrop-blur-3xl">
            <img src="/imgs/watch-1.webp" alt="Product" class="w-full h-full object-cover transition-transform duration-1000 group-hover:scale-105" />
            
            <!-- Overlay Info -->
            <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-700"></div>
            
            <div class="absolute bottom-8 left-8 right-8 flex items-end justify-between translate-y-4 group-hover:translate-y-0 transition-transform duration-700">
              <div class="space-y-1">
                <span class="text-[10px] font-mono text-white/40 uppercase tracking-[0.2em]">Mã định danh</span>
                <p class="text-xs font-mono text-white/80">SYMB-001-ALPHA</p>
              </div>
              <div class="text-right">
                <span class="text-[10px] font-mono text-white/40 uppercase tracking-[0.2em]">Giá trị</span>
                <p class="text-2xl font-spiky tracking-tighter">
                  <span x-text="(dropData?.price || 0).toLocaleString('vi-VN')"></span>
                  <span class="text-sm text-white/40 ml-1">VND</span>
                </p>
              </div>
            </div>
          </div>

          <!-- Status Badge -->
          <div class="absolute top-6 left-6">
            <div class="px-4 py-2 rounded-full backdrop-blur-xl border border-white/10 flex items-center gap-3"
                 :class="{
                   'bg-yellow-500/10 border-yellow-500/20': phase === 'WAITING',
                   'bg-green-500/10 border-green-500/20': phase === 'LIVE',
                   'bg-red-500/10 border-red-500/20': phase === 'SOLD_OUT' || phase === 'ENDED'
                 }">
              <span class="w-2 h-2 rounded-full animate-pulse"
                    :class="{
                      'bg-yellow-500': phase === 'WAITING',
                      'bg-green-500': phase === 'LIVE',
                      'bg-red-500': phase === 'SOLD_OUT' || phase === 'ENDED'
                    }"></span>
              <span class="text-[10px] font-bold uppercase tracking-[0.2em]"
                    :class="{
                      'text-yellow-500': phase === 'WAITING',
                      'text-green-500': phase === 'LIVE',
                      'text-red-500': phase === 'SOLD_OUT' || phase === 'ENDED'
                    }"
                    x-text="phase === 'WAITING' ? 'Sắp bắt đầu' : phase === 'LIVE' ? 'Đang diễn ra' : phase === 'SOLD_OUT' ? 'Đã bán hết' : 'Đã kết thúc'"></span>
            </div>
          </div>
        </div>
      </div>

      <!-- Product Info & Actions -->
      <div class="lg:col-span-5 space-y-12 reveal" style="transition-delay: 400ms">
        <div class="space-y-6">
          <div class="space-y-2">
            <h2 class="text-5xl md:text-7xl font-spiky tracking-tighter leading-none uppercase" x-text="dropData?.product_name"></h2>
            <p class="text-lg text-white/60 font-playfair italic">Cá thể cộng sinh đầu tiên của bộ sưu tập.</p>
          </div>

          <div class="flex items-center gap-12 py-8 border-y border-white/5">
            <div class="space-y-1">
              <span class="text-[10px] font-mono text-white/30 uppercase tracking-[0.3em]">Khả dụng</span>
              <div class="flex items-baseline gap-2">
                <span class="text-2xl font-spiky" x-text="dropData?.available"></span>
                <span class="text-xs text-white/40 uppercase tracking-widest">/ <span x-text="dropData?.drop_size"></span> Đơn vị</span>
              </div>
            </div>
            <div class="space-y-1">
              <span class="text-[10px] font-mono text-white/30 uppercase tracking-[0.3em]">Loại Drop</span>
              <p class="text-xs font-bold uppercase tracking-widest">Phát hành Alpha</p>
            </div>
          </div>
        </div>

        <!-- System Logs (FOMO) -->
        <div x-show="phase === 'LIVE'" class="space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="text-[10px] font-mono text-white/40 uppercase tracking-[0.3em]">Nhật ký Giao thức Trực tiếp</h3>
            <span class="text-[9px] font-mono text-green-500/60 uppercase animate-pulse">Đang đồng bộ...</span>
          </div>
          <div class="p-6 rounded-2xl bg-white/[0.02] border border-white/5 backdrop-blur-sm space-y-3">
            <template x-for="status in fomoStatuses" :key="status.phone">
              <div class="flex items-center justify-between text-[10px] font-mono">
                <div class="flex items-center gap-3">
                  <span class="text-white/20">>></span>
                  <span class="text-white/60" x-text="status.phone"></span>
                </div>
                <span class="text-white/40 italic" x-text="status.action"></span>
                <span class="text-white/20" x-text="status.time"></span>
              </div>
            </template>
          </div>
        </div>

        <!-- Action -->
        <div class="space-y-6">
          <div class="relative group" @mousemove="handleMagnetic" @mouseleave="resetMagnetic">
            <button @click="openModal()" 
                    :disabled="isDisabled" 
                    class="w-full py-6 rounded-full text-xs font-bold uppercase tracking-[0.4em] transition-all duration-500 overflow-hidden relative group"
                    :class="!isDisabled ? 'bg-white text-black hover:scale-[1.02]' : 'bg-white/5 text-white/20 cursor-not-allowed'"
                    :style="magneticStyle">
              <span class="relative z-10" x-text="phase === 'WAITING' ? 'Chờ bắt đầu' : !isDisabled ? 'Tham Gia Đấu Trường' : 'Đã kết thúc'"></span>
              
              <!-- Hover Glow -->
              <div x-show="!isDisabled" class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-full group-hover:animate-shimmer"></div>
            </button>
          </div>
          <p class="text-[9px] text-center text-white/30 uppercase tracking-[0.2em]">
            Luật Giao thức: Một slot cho mỗi định danh. Cửa sổ thanh toán 60 giây.
          </p>
        </div>
      </div>
    </div>
  `;
}
