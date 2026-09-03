import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";
import path from "node:path";

export default defineConfig({
    plugins: [
        // tailwindcss(),
    ],
    publicDir: false,
    input: {
        main: path.resolve(process.cwd(), "view/scripts/main.js"),
        // style: path.resolve(process.cwd(), "view/css/main.css"),
    },
    build: {
        outDir: path.resolve(process.cwd(), "public/assets"),
        emptyOutDir: false,
        sourcemap: true,
        // manifest: true,
        rolldownOptions: {
            output: {
                entryFileNames: "[name].js",
                assetFileNames: "[name][extname]",
                // chunkFileNames: "[name][extname]",
            },
        },
    },
});
