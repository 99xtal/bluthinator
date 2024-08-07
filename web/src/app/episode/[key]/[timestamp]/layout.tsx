import { ReactNode } from "react"

export default function Layout({ frame, frame_selector, children }: { frame: ReactNode, frame_selector: ReactNode, children: ReactNode }) {
    return (
        <div className="flex flex-col gap-8">
            {frame}
            {frame_selector}
        </div>
    )
}