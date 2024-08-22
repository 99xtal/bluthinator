"use client";

import { Button } from "../elements";

export default function ShareToFacebook({ url, hashtag }: { url: string, hashtag?: string }) {
    const handleShareToFacebook = () => {
        FB.ui({
          method: 'share',
          hashtag: hashtag,
          href: url,
        }, function(response: any){});
      }

    return (
        <Button onClick={handleShareToFacebook}>Share To Facebook</Button>
    )
}