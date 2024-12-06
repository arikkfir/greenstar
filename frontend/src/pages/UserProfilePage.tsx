import {useTheme} from "@mui/material";
import {UserProfile} from "@descope/react-sdk";

export function UserProfilePage() {
    const theme = useTheme()
    return (
        <UserProfile widgetId="user-profile-widget" theme={theme.palette.mode}/>
    )
}
