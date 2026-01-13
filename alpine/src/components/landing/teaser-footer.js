// Teaser Footer Web Component
class TeaserFooter extends HTMLElement {
  constructor() {
    super();
    // Use light DOM
  }

  connectedCallback() {
    this.render();
  }

  render() {
    this.innerHTML = /*html*/ `
            <footer class="py-12 px-6 sm:px-8 md:px-12 border-t theme-border bg-gradient-to-t from-black/60 via-black/40 to-transparent backdrop-blur-xl">
        <div class="max-w-7xl mx-auto space-y-8">
            <div class="flex flex-col md:flex-row items-center justify-between gap-8">
                <!-- Brand & Status -->
                <div class="flex flex-col items-center md:items-start gap-2">
                    <p class="text-lg font-black tracking-tighter text-white transition-all duration-500 hover:tracking-widest cursor-default">
                        DONALD VIBE LABS.
                    </p>
                    <div class="flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 border border-primary/20">
                        <span class="relative flex h-2 w-2">
                            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
                            <span class="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
                        </span>
                        <span class="text-[10px] font-mono text-primary uppercase tracking-widest">Hệ thống đang vận hành</span>
                    </div>
                </div>

                <!-- Links -->
                <div class="flex flex-wrap justify-center items-center gap-x-8 gap-y-4">
                    <a
                        href="mailto:studio@donaldvibe.xyz"
                        class="group text-sm text-muted-foreground hover:text-white transition-all duration-300 flex items-center gap-2"
                    >
                        <span class="w-1 h-1 bg-primary rounded-full opacity-0 group-hover:opacity-100 transition-opacity"></span>
                        studio@donaldvibe.xyz
                    </a>
                    <a
                        href="#verify"
                        @click.prevent="$dispatch('route', { page: 'verify' })"
                        class="group text-sm text-muted-foreground hover:text-white transition-all duration-300 flex items-center gap-2"
                    >
                        <span class="w-1 h-1 bg-primary rounded-full opacity-0 group-hover:opacity-100 transition-opacity"></span>
                        Xác thực SYMBICODE
                    </a>
                    <a
                        href="#policy"
                        @click.prevent="$dispatch('route', { page: 'policy' })"
                        class="group text-sm text-muted-foreground hover:text-white transition-all duration-300 flex items-center gap-2"
                    >
                        <span class="w-1 h-1 bg-primary rounded-full opacity-0 group-hover:opacity-100 transition-opacity"></span>
                        Quy trình Bảo mật
                    </a>
                    <a
                        href="#terms"
                        @click.prevent="$dispatch('route', { page: 'terms' })"
                        class="group text-sm text-muted-foreground hover:text-white transition-all duration-300 flex items-center gap-2"
                    >
                        <span class="w-1 h-1 bg-primary rounded-full opacity-0 group-hover:opacity-100 transition-opacity"></span>
                        Điều khoản Giao thức
                    </a>
                </div>
            </div>

            <!-- Bottom Bar -->
            <div class="pt-8 border-t border-white/5 flex flex-col sm:flex-row items-center justify-between gap-4">
                <p class="text-[10px] font-mono text-zinc-600 uppercase tracking-[0.3em]">
                    © 2025 GIAO THỨC CỘNG SINH. BẢO LƯU MỌI QUYỀN.
                </p>
                <div class="flex items-center gap-4">
                    <div class="h-px w-8 bg-white/10"></div>
                    <span class="text-[10px] font-mono text-zinc-600 uppercase tracking-[0.3em]">v1.0.4-ổn định</span>
                </div>
            </div>
        </div>
    </footer>
`;
  }
}

customElements.define("teaser-footer", TeaserFooter);
