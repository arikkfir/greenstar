import { fileURLToPath } from "url"
import path, { dirname } from "path"
import fs from "fs"

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

const defs = fs
    .readdirSync(__dirname)
    .filter((file) => file.endsWith(".graphql"))
    .map((file) => fs.readFileSync(path.join(__dirname, file), "utf8"))
    .join("\n")

export const TypeDefinitions = defs
