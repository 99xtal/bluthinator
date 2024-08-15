import Image from "next/image";

import { defonteRegular } from "~/fonts";
import GoBackLink from "~/app/components/GoBackLink";

export default async function Page({ params }: { params: { key: string, timestamp: string, caption: string } }) {
    const url = `${process.env.NEXT_PUBLIC_API_HOST}/caption/${params.key}/${params.timestamp}?b=${params.caption}`;
 
    return (
      <div>
        <div className="flex justify-start">
          <GoBackLink className={`${defonteRegular.className} text-md text-theme-black hover:underline`}>
            Back to Caption
          </GoBackLink>
        </div>
        <div className="flex justify-center items-center">
          <Image src={url} alt={params.caption} width={640} height={360} className="outline outline-4 outline-theme-black" />
        </div>
      </div>
    );
}