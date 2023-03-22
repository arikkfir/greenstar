export class User {
    id?: string
    email?: string
    name?: string
    given_name?: string
    family_name?: string
    picture?: string

    constructor(json: any) {
        this.id = json.id
        this.email = json.email
        this.name = json.name
        this.given_name = json.given_name
        this.family_name = json.family_name
        this.picture = json.picture
    }
}

export type LoginStatus = User | "pending" | "forbidden" | "error"
