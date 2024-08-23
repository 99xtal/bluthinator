"use client";

import Image from "next/image";

import { Button } from "../elements";
import logo from "~/assets/img/facebook_logo_secondary.png";

export default function ShareToFacebook({ url, hashtag }: { url: string, hashtag?: string }) {
    const handleShareToFacebook = () => {
        FB.ui({
          method: 'share',
          hashtag: hashtag,
          href: url,
        }, function(response: any){});
      }

    return (
        <Button onClick={handleShareToFacebook} className="px-6 py-2">
          <Image 
            src={logo}
            alt="Share to Facebook"
            width={32}
            height={32}
          />
        </Button>
    )
}