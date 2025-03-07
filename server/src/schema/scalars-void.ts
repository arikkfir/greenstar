import { GraphQLError, GraphQLScalarType } from "graphql"

export const VoidScalar = new GraphQLScalarType({
    name: "Void",
    description: "Represents an empty response",
    serialize() {
        return null
    },
    parseValue() {
        throw new GraphQLError("Void cannot be used as input")
    },
    parseLiteral() {
        throw new GraphQLError("Void cannot be used as input")
    },
})
