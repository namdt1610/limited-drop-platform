echo "=== ALPINE (npm deps) ==="
echo "Direct dependencies:"
grep '"alpinejs"\|"tailwindcss"\|"franken-ui"' alpine/package.json
echo ""
echo "DevDependencies:"
grep '"vite"\|"@vitejs/plugin"\|"bun"\|"tailwindcss"' alpine/package.json | grep devDependencies -A 10
echo ""

echo "=== BACKEND (Go deps) ==="
echo "Direct Go imports used in source files:"
find backend/internal backend/cmd -name "*.go" -type f | xargs grep -h "^import\|^\s*\"" | sort | uniq | grep -v "^--"
