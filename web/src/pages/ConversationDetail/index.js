import React, { useContext, useEffect, useState, useCallback } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { API, showError, showSuccess } from '../../helpers';
import { Chat, Button } from '@douyinfe/semi-ui';
import {
  IconArrowLeft,
  IconSidebar,
  IconChevronDown,
} from '@douyinfe/semi-icons';
import './index.scss';
const roleInfo = {
  user: {
    name: '用户',
    avatar:
      'https://ziyile-1347238265.cos.ap-guangzhou.myqcloud.com/Beijixing/assets/star.png',
  },
  assistant: {
    name: '小星星',
    avatar:
      'https://ziyile-1347238265.cos.ap-guangzhou.myqcloud.com/Beijixing/assets/icon.png',
  },
  system: {
    name: '小星星',
    avatar:
      'https://ziyile-1347238265.cos.ap-guangzhou.myqcloud.com/Beijixing/assets/icon.png',
  },
};

const ConversationDetail = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { id, type } = location.state || {};
  const [message, setMessage] = useState();

  useEffect(() => {
    loadConversation();
  }, []);

  const detailStr = (str) => {
    const match = str.match(/<think>([\s\S]*?)<\/think>([\s\S]*)/);
    if (match) {
      const contentInside = match[1];
      const contentOutside = match[2];
      const formattedReasoningContent = contentInside
        .split('\n')
        .map((line) => `> ${line}`)
        .join('\n');
      return formattedReasoningContent + contentOutside;
    } else {
      return str;
    }
  };

  const loadConversation = async () => {
    const res = await API.get(`/api/conversation/${id}/messages`);
    const { success, message, data } = res.data;
    if (success) {
      setMessage(data);

      setMessage((messages) => {
        return messages.map((element) => {
          const newElement = { ...element };

          try {
            if (newElement.content_type === 'text') {
              newElement.content = detailStr(newElement.content);
            } else if (newElement.content_type === 'image') {
              const parsed = JSON.parse(newElement.content);
              if (Array.isArray(parsed.data)) {
                newElement.content = parsed.data.map((e) => ({
                  type: 'image_url',
                  image_url: { url: e.url },
                }));
              } else {
                console.warn('Invalid image data format', parsed);
                newElement.content = [];
              }
            } else {
              newElement.content = JSON.parse(newElement.content);
            }
          } catch (error) {
            console.error('Error parsing content:', error, newElement);
            // 这里可以考虑是否保留原始 content 或设置为 null
          }

          return newElement;
        });
      });
    } else {
      showError(message);
    }
    // setLoading(false);
  };

  const commonOuterStyle = {
    border: '1px solid var(--semi-color-border)',
    borderRadius: '16px',
    margin: '8px 16px',
  };
  //输入框为空
  const renderInputArea = useCallback((props) => {
    return <></>;
  }, []);

  return (
    <div className='chatcontent'>
      <div style={{ width: '100%', display: 'flex' }}>
        <Button
          icon={<IconArrowLeft />}
          theme='solid'
          style={{ marginRight: 10 }}
          onClick={() => {
            navigate('/conversation');
          }}
        >
          返回
        </Button>
      </div>

      <Chat
        mode={'userBubble'}
        renderInputArea={renderInputArea}
        style={commonOuterStyle}
        chats={message}
        roleConfig={roleInfo}
      />
    </div>
  );
};

export default ConversationDetail;
