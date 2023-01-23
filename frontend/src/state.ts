import {Atom, atom} from "jotai";

export class User {
    constructor(
        public readonly id: string,
        public readonly firstName: string,
        public readonly lastName: string,
        public readonly jwt: string,
    ) {
    }

    get name(): string {
        return `${this.firstName} ${this.lastName}`;
    }
}

let userAtom: Atom<User>;
const stringValue = localStorage.getItem('user');
if (!stringValue) {
    window.location.href = "http://localhost:8000/auth/google"
    userAtom = atom<User>(new User("ANONYMOUS", "Anonymous", "Anonymous", "INVALID"));
} else {
    const objectValue = JSON.parse(stringValue);
    userAtom = atom<User>(new User(objectValue.id, objectValue.firstName, objectValue.lastName, objectValue.jwt));
}

const state = {
    user: userAtom,
}

export default state
