import React, { useEffect, useState } from 'react';
import { Tooltip } from 'antd';
import { GetCheckCodeImage } from '@/services/ant-design-pro/api';


const xInitImgData='/check_code_refresh.png';

const CheckCodeImage: React.FC = () => {

    const[xImgData,setImgData]=useState<string>(xInitImgData);

    const loadImgData=async ()=>{

        const xRemoteImgData=await GetCheckCodeImage();

        if(xRemoteImgData){
            setImgData(xRemoteImgData);
        }
        
    };

    useEffect(()=>{
        loadImgData();
    },[]);

    return(
      <Tooltip title="点击刷新验证码">
            <img style={{cursor:'pointer'}} height={35} width={88} src={xImgData} onClick={loadImgData} />
      </Tooltip>
    );
};


export default CheckCodeImage;