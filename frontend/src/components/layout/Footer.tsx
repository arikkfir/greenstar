import {Link} from "react-router-dom";
import {Box} from "@mui/material";
import pkg from "../../../package.json";

export function Footer() {
    return (
        <Box sx={{
            display: "flex",
            flexDirection: "row",
            justifyContent: "center",
            alignContent: "center",
        }}>
            Copyright © 2017 - {new Date().getFullYear()}&nbsp;&nbsp;
            <Link to="https://github.com/arikkfir">Arik Kfir</Link>&nbsp;&nbsp;
            All Rights Reserved ({pkg.version})
        </Box>
    )
}
