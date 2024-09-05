import { NextRequest, NextResponse } from "next/server";

export const revalidate = 0;

export async function GET(req: NextRequest) {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/random`, { cache: "no-store" });
    const data = await response.json();
  
    const path = `/episode/${data.episode}/${data.timestamp}`;
    return NextResponse.redirect(new URL(path, req.url));
}