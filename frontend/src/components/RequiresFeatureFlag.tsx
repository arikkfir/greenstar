import {PropsWithChildren} from "react";
import {useFeatureIsOn} from "@growthbook/growthbook-react";

interface RequiresFeatureFlagProps {
    mustBeOn: string[]
}

export function RequiresFeatureFlag({mustBeOn, children}: PropsWithChildren<RequiresFeatureFlagProps>) {
    if (mustBeOn.every(ff => useFeatureIsOn(ff))) {
        return children
    } else {
        return <></>
    }
}
