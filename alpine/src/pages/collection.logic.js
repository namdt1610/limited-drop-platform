export function CollectionLogic() {
  return {
    init() {
      this.initRevealObserver();
    },
    initRevealObserver() {
      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              entry.target.classList.add("reveal-visible");
            }
          });
        },
        { threshold: 0.1 }
      );
      setTimeout(() => {
        document
          .querySelectorAll(".reveal")
          .forEach((el) => observer.observe(el));
      }, 500);
    },
    goHome() {
      window.location.href = "/";
    },
  };
}
