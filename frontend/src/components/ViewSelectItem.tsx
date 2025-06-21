type ViewProps = {
  name: string;
  creator: string;
  issue_num: number;
  destination: string;
};

export default function ViewSelectItem(props: ViewProps) {
  const handleClick = () => {
    console.log(props.destination);
  };
  return (
    <div
      className="flex flex-row align-center w-4/5 p-4 m-2 rounded-lg bg-gray-800 hover:bg-gray-700 cursor-pointer"
      onClick={handleClick}
    >
      <span className="inline-block w-1/2 font-bold text-xl p-2 overflow-hidden text-ellipsis whitespace-nowrap">
        {props.name}
      </span>
      <span className="inline-block w-3/8 font-normal text-lg p-2 overflow-hidden text-ellipsis whitespace-nowrap">
        {props.creator}
      </span>
      <span className="inline-block w-1/8 font-normal text-lg p-2">
        {props.issue_num}
      </span>
    </div>
  );
}
