import React, { useEffect, useState } from 'react';
import moment from 'moment';
import 'moment/locale/zh-cn';
import {
  CtlIngressActionBackendMap,
  CtlIngressActionStaticMap,
  CtlIngressSiteRuleInfo,
  DataMapAppDataIngressActionType,
  DataMapHttpMethod,
  DataMapHttpTargetItem,
  DataMapRuleStringOp,
} from '@/services/ant-design-pro/api';
import { Modal, Spin } from 'antd';
import ProDescriptions from '@ant-design/pro-descriptions';

export type InfoRuleFormProps = {
  onModalVisible: (flag: boolean) => void;
  modalVisible: boolean;
  namespace_id?: string;
  id?: string;
};

const InfoRuleForm: React.FC<InfoRuleFormProps> = (props) => {
  const [xIsLoading, setIsLoading] = useState<boolean>(true);
  const [xIngressActionValueMapData, setIngressActionValueMapData] = useState<any>();

  const [xInfoData, setInfoData] = useState<API.AppDataIngressSiteRule>();

  const onFormVisible = (flag: boolean) => {
    moment.locale('zh-cn');
    return props.onModalVisible(flag);
  };

  const loadIngressActionValueMapData = async (infoData: API.AppDataIngressSiteRule) => {
    const xDataList: any[] = [];

    const xActionType = infoData.action_type;

    const xNamespaceId = props.namespace_id;
    const xReq = { data_id: xNamespaceId };

    if (xActionType == 'backend') {
      const xRespData = await CtlIngressActionBackendMap(xReq);
      if (xRespData && xRespData.data) {
        xRespData.data.forEach((xItem) => {
          xDataList.push({ label: xItem.data_name, value: xItem.data_id });
        });
      }
    }

    if (xActionType == 'static') {
      const xRespData = await CtlIngressActionStaticMap(xReq);
      if (xRespData && xRespData.data) {
        xRespData.data.forEach((xItem) => {
          xDataList.push({ label: xItem.data_name, value: xItem.data_id });
        });
      }
    }

    setIngressActionValueMapData(xDataList);
  };

  const loadDataInfo = async (id?: string) => {
    setIsLoading(true);

    if (id) {
      const respData = await CtlIngressSiteRuleInfo({ data_id: id });
      if (respData && respData.error_code === '0') {
        if (respData.data) {
          const xDataInfo = respData.data;
          await loadIngressActionValueMapData(xDataInfo);
          setInfoData(xDataInfo);
        }
      }
    }

    setIsLoading(false);
  };

  // 加载
  useEffect(() => {
    if (props.modalVisible) {
      loadDataInfo(props.id);
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [props.id, props.modalVisible]);

  return (
    <Modal
      title={'路由规则详情'}
      width="80%"
      visible={props.modalVisible}
      destroyOnClose={true}
      onCancel={() => {
        onFormVisible(false);
      }}
      footer={null}
      key="view-modal-form"
    >
      <Spin spinning={xIsLoading} delay={100}>
        <ProDescriptions
          dataSource={xInfoData}
          bordered
          column={2}
          columns={[
            {
              title: '系统ID',
              dataIndex: 'data_id',
              hideInSearch: true,
              hideInTable: true,
              hideInDescriptions: true,
            },

            {
              title: '规则序号',
              dataIndex: 'order_no',
              hideInSearch: true,
            },

            {
              title: '规则名称',
              dataIndex: 'title',
              hideInSearch: true,
              ellipsis: true,
            },

            {
              title: 'HTTP方法',
              dataIndex: 'http_method',
              valueType: 'select',
              hideInSearch: true,
              valueEnum: DataMapHttpMethod,
            },

            {
              title: '规则匹配目标',
              dataIndex: 'match_target',
              valueType: 'select',
              hideInSearch: true,
              valueEnum: DataMapHttpTargetItem,
            },

            {
              title: '规则匹配方式',
              dataIndex: 'match_op',
              valueType: 'select',
              hideInSearch: true,
              valueEnum: DataMapRuleStringOp,
            },
            {
              title: '规则匹配内容',
              dataIndex: 'match_value',
              hideInSearch: true,
              ellipsis: true,
            },

            {
              title: '响应类型',
              dataIndex: 'action_type',
              valueType: 'select',
              hideInSearch: true,
              valueEnum: DataMapAppDataIngressActionType,
            },

            {
              title: '响应内容',
              dataIndex: 'action_value',
              valueType: 'select',
              hideInSearch: true,
              fieldProps: { options: xIngressActionValueMapData },
            },
          ]}
        />
      </Spin>
    </Modal>
  );
};

export default InfoRuleForm;
