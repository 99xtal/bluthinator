import { NextRequest, NextResponse } from "next/server";

export const revalidate = 0;

export async function GET(req: NextRequest) {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/random`, { cache: "no-store" });
    const data = await response.json();
  
    const url = req.nextUrl;
    url.pathname = `/episode/${data.episode}/${data.timestamp}`;
    return NextResponse.redirect(url);
}