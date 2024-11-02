import {UserProfile} from "@descope/react-sdk";
import {useTheme} from "@mui/material";

export function UserProfilePage() {
    const theme = useTheme()
    return (
        <UserProfile widgetId="user-profile-widget" theme={theme.palette.mode}/>
    )
}
