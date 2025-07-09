import StatusIcon from "./StatusIcon";
import PriorityIcon from "./PriorityIcon";
import AssignedIcon from "./AssignedIcon";
import DateFormat from "./DateFormat";

type SelectableItemProps = {
  name: string;
  status: string;
  priority: string;
  assigned: string;
  startDate: string;
  endDate: string;
  destination: string;
};

export default function IssueSelectItem(props: SelectableItemProps) {
  const handleClick = () => {
    if (props.destination != "") {
      console.log(props.destination);
    }
  };
  return (
    <div
      className="flex flex-row items-center w-4/5 px-2 py-1 my-2 text-gray-200 rounded-lg bg-gray-800 hover:bg-gray-700 active:bg-gray-600 cursor-pointer transition duration-200"
      onClick={handleClick}
    >
      <span className="inline-block w-7/10 font-semibold text-base p-2 overflow-hidden text-ellipsis whitespace-nowrap">
        {props.name}
      </span>
      <span className="w-1/20 font-bold text-lg p-2">
        {StatusIcon(props.status)}
      </span>
      <span className="w-1/20 font-bold text-lg p-2">
        {PriorityIcon(props.priority)}
      </span>
      <span className="w-1/20 font-bold text-lg p-2">
        {AssignedIcon(props.assigned)}
      </span>
      <span className="w-2/5 font-normal text-sm p-2 pl-30">
        {DateFormat(props.startDate)} - {DateFormat(props.endDate)}
      </span>
    </div>
  );
}
