import "./Layout.scss"
import { buildErrorExtra, ErrorCard } from "../../components/ErrorCard.tsx"

export interface ErrorPageProps {
    title: string
    subTitle?: string
    error?: Error
}

export function ErrorPage({ title, subTitle, error }: ErrorPageProps) {
    return (
        <main className="error-page">
            <ErrorCard title={title}
                       subTitle={subTitle}
                       text={error?.message}
                       extra={error && buildErrorExtra(error)} />
        </main>
    )
}
