export function HeaderComponent() {
  return /*html*/ `
    <header class="reveal py-8 px-6 md:px-12">
      <div class="max-w-screen-2xl mx-auto flex items-center justify-between">
        <div class="flex items-center gap-6">
          <a href="#landing" @click.prevent="$dispatch('route', 'landing')" class="group">
            <img src="/imgs/logo-png.webp" alt="DONALD" class="w-20 md:w-24 h-auto object-contain transition-transform duration-500 group-hover:scale-110" />
          </a>
          <div class="hidden md:block h-8 w-[1px] bg-white/10"></div>
          <div class="hidden md:block">
            <h1 class="text-[10px] font-mono tracking-[0.4em] text-white/40 uppercase mb-1" x-text="dropData?.name"></h1>
            <p class="text-xs font-medium tracking-widest uppercase" x-text="dropData?.product_name"></p>
          </div>
        </div>

        <!-- System Status -->
        <div class="flex items-center gap-8">
          <div class="hidden lg:flex flex-col items-end">
            <span class="text-[9px] font-mono text-white/30 uppercase tracking-[0.3em] mb-1">Trạng thái Hệ thống</span>
            <div class="flex items-center gap-2">
              <span class="w-1.5 h-1.5 rounded-full bg-green-500 animate-pulse"></span>
              <span class="text-[10px] font-mono text-green-500/80 uppercase tracking-widest">Giao thức Kích hoạt</span>
            </div>
          </div>
          
          <div class="flex flex-col items-end">
            <span class="text-[9px] font-mono text-white/30 uppercase tracking-[0.3em] mb-1" x-text="phase === 'WAITING' ? 'Đếm ngược' : 'Thời gian còn lại'"></span>
            <div class="text-xl md:text-2xl font-spiky tracking-tighter" x-text="countdown"></div>
          </div>
        </div>
      </div>
    </header>
  `;
}
