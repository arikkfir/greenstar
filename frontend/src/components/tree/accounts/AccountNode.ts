import {Account} from "../../../client/account.ts";

export interface AccountNode extends Account {
    children: Account[]
}
