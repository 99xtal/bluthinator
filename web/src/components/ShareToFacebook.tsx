"use client";

import Image from "next/image";

import { Button } from "../elements";
import logo from "~/assets/img/facebook_logo_secondary.png";
import { usePathname } from "next/navigation";
import { logEvent } from "~/utils/firebase";

export default function ShareToFacebook({ hashtag }: { hashtag?: string }) {
  const pathname = usePathname();

  const handleShareToFacebook = () => {
      logEvent("share_meme_facebook", { pathname });
      FB.ui({
        method: 'share',
        hashtag: hashtag,
        href: `https://bluthinator.com${pathname}`,
      }, function(response: any){});
    }

  return (
      <Button onClick={handleShareToFacebook} title="Share To Facebook" className="px-6 py-2">
        <Image 
          src={logo}
          alt="Share to Facebook"
          width={32}
          height={32}
        />
      </Button>
  )
}