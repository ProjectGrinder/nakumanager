import IssueCount from "../IssueCount";

type ViewProps = {
  name: string;
  creator: string;
  issue_list: string[];
  destination: string;
};

export default function ViewSelectItem(props: ViewProps) {
  const handleClick = () => {
    console.log(props.destination);
  };
  return (
    <div
      className="flex flex-row items-center w-4/5 p-2 my-2 text-gray-200 rounded-lg bg-gray-800 hover:bg-gray-700 active:bg-gray-600 cursor-pointer transition duration-200"
      onClick={handleClick}
    >
      <span className="inline-block w-1/2 font-bold text-lg p-2 overflow-hidden text-ellipsis whitespace-nowrap">
        {props.name}
      </span>
      <span className="inline-block w-2/5 font-normal text-base p-2 overflow-hidden text-ellipsis whitespace-nowrap">
        {props.creator}
      </span>
      <span className="inline-block w-1/10 font-normal text-base p-2">
        {IssueCount(props.issue_list)}
      </span>
    </div>
  );
}
