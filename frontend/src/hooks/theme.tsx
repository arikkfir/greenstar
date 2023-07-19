import {ListItemButtonProps, PaletteMode, useMediaQuery} from "@mui/material";
import {createTheme} from "@mui/material/styles";
import React, {useState} from "react";
import {Link as RouterLink, LinkProps as RouterLinkProps} from 'react-router-dom';

const MuiListItemButtonBehavior = React.forwardRef<
    HTMLAnchorElement,
    Omit<RouterLinkProps, 'to'> & { href: RouterLinkProps['to'] }
>((props, ref) => {
    const {href, ...other} = props;
    return <RouterLink ref={ref} to={href} {...other} />;
});

export function useTheme() {
    const osPaletteMode = useMediaQuery('(prefers-color-scheme: dark)') ? 'dark' : 'light';
    const [preferredPaletteMode] = useState<PaletteMode | null>(null)
    return React.useMemo(
        () => createTheme({
            components: {
                MuiListItemButton: {
                    defaultProps: {
                        component: MuiListItemButtonBehavior,
                    } as ListItemButtonProps
                }
            },
            palette: {
                mode: preferredPaletteMode || osPaletteMode,
            },
        }),
        [osPaletteMode, preferredPaletteMode],
    )
}
