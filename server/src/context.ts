import { BaseContext } from "@apollo/server"
import { DataLayer } from "./data/data_layer.js"

export interface Context extends BaseContext {
    data: DataLayer
}
