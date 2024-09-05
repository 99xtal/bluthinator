import Image from "next/image";

import { ffBlurProMedium } from "~/fonts";
import ShareToFacebook from "~/app/components/ShareToFacebook";
import Link from "next/link";
import CopyLinkButton from "~/app/components/CopyLinkButton";
import { Metadata } from "next";

type Props = {
  params: { key: string, timestamp: string, caption: string }
}

export async function generateMetadata(
  { params }: Props,
): Promise<Metadata> {
  return {
    openGraph: {
      title: "Bluthinator",
      description: "Check out this meme from Bluthinator!",
      images: [
        {
          url: `${process.env.NEXT_PUBLIC_API_HOST}/caption/${params.key}/${params.timestamp}?b=${params.caption}`,
          width: 720,
          height: 405,
          alt: "Bluthinator meme",
        },
      ]
    }
  }
}

export default async function Page({ params }: { params: { key: string, timestamp: string, caption: string } }) {
    const url = `${process.env.NEXT_PUBLIC_API_HOST}/caption/${params.key}/${params.timestamp}?b=${params.caption}`;
 
    return (
      <div className="flex justify-center items-center">
        <div className="flex flex-col gap-4">
          <Link href={`/img/caption/${params.key}/${params.timestamp}/${params.caption}`}>
            <Image src={url} alt={params.caption} width={640} height={360} className="outline outline-4 outline-theme-black" />
          </Link>
          <div className="flex justify-between items-center">
            <Link href={`/episode/${params.key}/${params.timestamp}`} className={`${ffBlurProMedium.className} text-lg text-theme-black hover:underline`}>
              &larr; Back to Episode
            </Link>
            <div className="flex gap-2 items-center">
              <CopyLinkButton />
              <ShareToFacebook hashtag="#bluthinator" />
            </div>
          </div>
        </div>
      </div>
    );
}