import { redirect } from "next/navigation";

const getRandomFrame = async () => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/random`, { cache: "no-store" });
    return response.json();
}

export default async function Page() {
    const randomFrame = await getRandomFrame();

    if (randomFrame) {
        redirect(`/episode/${randomFrame.episode}/${randomFrame.timestamp}`);
    }

    return null;
}