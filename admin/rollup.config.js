import { terser } from "rollup-plugin-terser";
import babel from "@rollup/plugin-babel";
import commonjs from "@rollup/plugin-commonjs";
import postcss from "rollup-plugin-postcss";
import resolve from "@rollup/plugin-node-resolve";
import replace from "@rollup/plugin-replace";
import dotenv from "dotenv";

require("dotenv").config();

export default {
    input: "src/index.js",
    output: {
        file: "public/bundle.js",
        format: "iife",
        sourcemap: true,
    },
    plugins: [
        resolve(),
        replace({
            preventAssignment: true,
            "process.env.NODE_ENV": JSON.stringify("production"),
            "process.env.BACKEND": JSON.stringify(process.env.BACKEND),
        }),
        postcss(),
        babel({
            exclude: "node_modules/**",
            presets: ["@babel/env", "@babel/preset-react"],
        }),
        commonjs(),
        terser(),
    ],
};
