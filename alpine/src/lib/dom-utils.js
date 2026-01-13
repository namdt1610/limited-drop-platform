/**
 * DOM & Animation Utilities - Shared across pages
 */

/**
 * Initialize Intersection Observer for reveal animations
 * Adds 'reveal-visible' class when element enters viewport
 */
export function initRevealObserver() {
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

  // Wait for DOM to be ready
  setTimeout(() => {
    document.querySelectorAll(".reveal").forEach((el) => observer.observe(el));
  }, 500);
}

/**
 * Handle magnetic cursor effect on element hover
 * Scales and applies transform based on mouse position
 */
export function handleMagnetic(element, event) {
  const rect = element.getBoundingClientRect();
  const centerX = rect.width / 2;
  const centerY = rect.height / 2;
  const x = event.clientX - rect.left - centerX;
  const y = event.clientY - rect.top - centerY;

  const distance = Math.sqrt(x * x + y * y);
  const maxDistance = Math.sqrt(centerX * centerX + centerY * centerY);
  const scale = Math.max(1, 1 + (maxDistance - distance) / (maxDistance * 2));

  return {
    transform: `translate(${x * 0.2}px, ${y * 0.2}px) scale(${scale})`,
  };
}

/**
 * Reset magnetic effect style
 */
export function resetMagneticStyle() {
  return { transform: "" };
}

/**
 * Format milliseconds to HH:MM:SS format
 */
export function formatTime(ms) {
  const seconds = Math.floor((ms / 1000) % 60);
  const minutes = Math.floor((ms / (1000 * 60)) % 60);
  const hours = Math.floor((ms / (1000 * 60 * 60)) % 24);
  return `${hours.toString().padStart(2, "0")}:${minutes
    .toString()
    .padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
}

/**
 * Format countdown with days, hours, minutes, seconds
 */
export function formatCountdown(ms) {
  const d = Math.floor(ms / (1000 * 60 * 60 * 24));
  const h = Math.floor((ms % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
  const m = Math.floor((ms % (1000 * 60 * 60)) / (1000 * 60));
  const s = Math.floor((ms % (1000 * 60)) / 1000);

  return {
    d: d.toString().padStart(2, "0"),
    h: h.toString().padStart(2, "0"),
    m: m.toString().padStart(2, "0"),
    s: s.toString().padStart(2, "0"),
  };
}
