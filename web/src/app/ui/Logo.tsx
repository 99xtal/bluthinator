import { defonteRegular } from "../../fonts";

export default function Logo() {
    return (
        <div className="transform -rotate-5">
            <h1 className={`${defonteRegular.className} text-3xl`}>bluthinator</h1>
        </div>
    );
}