import { message } from 'antd';
import { isNumber } from 'lodash';
import moment from 'moment';
import { getIntl } from 'umi';

export function GetFormatMessage(msg?: string) {

  let msgData =msg;


  if(msg&&msg?.indexOf('app.server.msg')>-1){
    const intl = getIntl();
    msgData = intl.formatMessage({ id: msg, defaultMessage: msg });
  }

  return msgData;
}

export function ShowInfoMessage(msg?: string) {
  const msgData = GetFormatMessage(msg);
  message.info(msgData);
}

export function ShowSuccessMessage(msg?: string) {
  const msgData = GetFormatMessage(msg);
  message.success(msgData);
}

export function ShowErrorMessage(msg?: string) {
  const msgData = GetFormatMessage(msg);

  message.error(msgData);
}

export function SetLocalCache(dataKey: string, dataVal: string) {
  window.localStorage.setItem(dataKey, dataVal);
}

export function GetLocalCache(dataKey: string): string {
  let xData = '';

  const xTempVal = window.localStorage.getItem(dataKey);

  if (xTempVal) {
    xData = xTempVal;
  }

  return xData;
}

export function RemoveLocalCache(dataKey: string) {
  window.localStorage.removeItem(dataKey);
}

export function ClearLocalCache() {
  window.localStorage.clear();
}

export function GetSelectMapKeyByVal(dataList: API.SelectItemNode[], val: any) {
  let xData = '';

  dataList.forEach((item) => {
    if (item.value == val) {
      xData = item.label;
    }
  });

  return xData;
}

export function ParseTimestamp(timeStamp?: number) {
  if (!timeStamp) {
    return 'N/A';
  }

  const xTime = timeStamp;

  const xUtcDate = moment.unix(xTime).toDate();

  return moment.parseZone(xUtcDate).format('YYYY-MM-DD HH:mm:ss');
}

export function ParseDeviceSize(data?: number) {
  let xData: string = 'N/A';

  if (data === undefined) {
    return xData;
  }

  let dataNum: number = data;

  const dataBase = 1024;

  if (dataNum < dataBase ** 1) {
    xData = `${dataNum.toFixed(4)}iB`;
    return xData;
  }

  if (dataNum < dataBase ** 2) {
    dataNum /= dataBase ** 1;
    xData = `${dataNum.toFixed(4)}KiB`;
    return xData;
  }

  if (dataNum < dataBase ** 3) {
    dataNum /= dataBase ** 2;
    xData = `${dataNum.toFixed(4)}MiB`;
    return xData;
  }

  if (dataNum < dataBase ** 4) {
    dataNum /= dataBase ** 3;
    xData = `${dataNum.toFixed(4)}GiB`;
    return xData;
  }

  if (dataNum < dataBase ** 5) {
    dataNum /= dataBase ** 4;
    xData = `${dataNum.toFixed(4)}TiB`;
    return xData;
  }

  if (dataNum < dataBase ** 6) {
    dataNum /= dataBase ** 5;
    xData = `${dataNum.toFixed(4)}PiB`;
    return xData;
  }

  if (dataNum < dataBase ** 7) {
    dataNum /= dataBase ** 6;
    xData = `${dataNum.toFixed(4)}EiB`;
    return xData;
  }

  return xData;
}

export function ParseDeviceTime(data?: number, unitSecond: boolean = false) {
  let xData: string = 'N/A';

  if (data === undefined) {
    return xData;
  }

  let dataNum: number = data;

  if (unitSecond) {
    dataNum = dataNum * 1000;
  }

  const dataBaseSec = 1000;
  const dataBaseMin = 1000 * 60;
  const dataBaseHour = 1000 * 60 * 60;
  const dataBaseDay = 1000 * 60 * 60 * 24;

  if (dataNum < dataBaseSec) {
    xData = `${dataNum.toFixed(4)}毫秒`;
    return xData;
  }

  if (dataNum < dataBaseMin) {
    dataNum /= dataBaseSec;
    xData = `${dataNum.toFixed(4)}秒`;
    return xData;
  }

  if (dataNum < dataBaseHour) {
    dataNum /= dataBaseMin;
    xData = `${dataNum.toFixed(4)}分`;
    return xData;
  }

  if (dataNum < dataBaseDay) {
    dataNum /= dataBaseHour;
    xData = `${dataNum.toFixed(4)}小时`;
    return xData;
  }

  dataNum /= dataBaseDay;
  xData = `${dataNum.toFixed(4)}天`;

  return xData;
}

export function GetDateTimeDiffSeconds(beginTime?: string,endTime?: string){

  let xData: string = 'N/A';


  if(!beginTime||!endTime){
    return xData;
  }

  const xDataVal=moment(endTime).diff(moment(beginTime),'seconds');

  if(isNumber(xDataVal)){

    if(xDataVal>=0){
      xData=`${xDataVal}秒`;
    }
  }

  return xData;

}