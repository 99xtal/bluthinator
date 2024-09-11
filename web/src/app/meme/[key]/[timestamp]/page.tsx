import { Metadata } from "next";
import Link from "next/link";
import { ffBlurProMedium } from "~/assets/fonts";
import { ShareToFacebook, CopyLinkButton } from "~/components";
import MemeResult from "~/components/MemeResult";

type Props = {
  params: { key: string, timestamp: string, caption: string }
}

export async function generateMetadata(
  { params }: Props,
): Promise<Metadata> {
  return {
    title: `Bluthinator | Meme "${Buffer.from(decodeURI(params.caption), 'base64').toString()}"`,
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

export default function Page({ params }: { params: { key: string, timestamp: string } }) {
    return (
      <div className="flex justify-center items-center">
        <div className="flex flex-col gap-4">
        <MemeResult episode={params.key} timestamp={params.timestamp} />
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