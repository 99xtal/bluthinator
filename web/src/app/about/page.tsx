import { ffBlurProMedium } from "~/fonts";
import Image from "next/image";
import { getFrameUrl } from "~/utils";

export default function Page() {
    return (
        <div className="w-full flex justify-center py-6">
            <div className="flex flex-col max-w-lg gap-8">
                <Image
                    src={getFrameUrl('S1E01', 89625, 'large')}
                    alt={"We Demand to be Taken Seriously"}
                    width={640}
                    height={360}
                    className="max-w-sm h-auto outline outline-4 outline-theme-black self-center"
                />
                <section>
                    <h3 className={`${ffBlurProMedium.className} text-2xl`}>About</h3>
                    <div className="flex flex-col gap-2">
                        <p>Episode frames from seasons 1-3 are indexed by subtitle text, so moments from the show are searchable by lines spoken by the characters.</p>
                        <p>This project is a labor of fandom, and is very much inspired by <a href="https://frinkiac.com" className="underline">Frinkiac</a>.</p>
                    </div>
                </section>
                <section>
                    <h3 className={`${ffBlurProMedium.className} text-2xl`}>Contact</h3>
                    <p>Your feedback and suggestions are very valuable! Feel free to email me at <a href="mailto:bluthinatorapp@gmail.com" className="underline">bluthinatorapp@gmail.com</a>.</p>
                </section>
                <section>
                    <h3 className={`${ffBlurProMedium.className} text-2xl`}>Credits</h3>
                    <ul className="list-disc">
                        <li>
                            <a href="https://twitter.com/reaperhulk" className="underline">Paul Kehrer</a>,&nbsp;
                            <a href="https://twitter.com/sirsean" className="underline">Sean Schulte</a>&nbsp;
                            &&nbsp;<a href="https://twitter.com/seriousallie" className="underline">Allie Young</a>, the creators of Frinkiac
                        </li>
                        <li><a href="https://ngallant.dev" className="underline">Nat Gallant</a> and Richard Rybarczyk</li>
                    </ul>
                </section>
            </div>
        </div>
    )
}